package service

import (
	"context"
	"errors"
	"fmt"

	"gogo/internal/cache"
	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/repository"

	"github.com/google/uuid"
)

// TerminalService handles terminal management business logic.
type TerminalService struct {
	terminalRepo  repository.TerminalRepository
	storeRepo     repository.StoreRepository
	heartbeatCache *cache.HeartbeatCache
	logRepo       repository.LogRepository
}

// NewTerminalService creates a new TerminalService.
func NewTerminalService(
	terminalRepo repository.TerminalRepository,
	storeRepo repository.StoreRepository,
	heartbeatCache *cache.HeartbeatCache,
	logRepo repository.LogRepository,
) *TerminalService {
	return &TerminalService{
		terminalRepo:  terminalRepo,
		storeRepo:     storeRepo,
		heartbeatCache: heartbeatCache,
		logRepo:       logRepo,
	}
}

// Create pre-registers a terminal with a generated UUID device_token.
func (s *TerminalService) Create(ctx context.Context, req dto.CreateTerminalReq) (*model.Terminal, error) {
	// Check store exists
	store, err := s.storeRepo.GetByID(ctx, req.StoreID)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	deviceToken := uuid.New().String()

	terminal := &model.Terminal{
		SN:          req.SN,
		Name:        req.Name,
		Type:        req.Type,
		StoreID:     req.StoreID,
		Status:      model.TerminalStatusOffline,
		DeviceToken: deviceToken,
	}

	if err := s.terminalRepo.Create(ctx, terminal); err != nil {
		return nil, fmt.Errorf("create terminal: %w", err)
	}

	return terminal, nil
}

// GetByID returns a terminal by ID.
func (s *TerminalService) GetByID(ctx context.Context, id int64) (*model.Terminal, error) {
	t, err := s.terminalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrTerminalNotFound
	}
	return t, nil
}

// GetBySN returns a terminal by SN.
func (s *TerminalService) GetBySN(ctx context.Context, sn string) (*model.Terminal, error) {
	t, err := s.terminalRepo.GetBySN(ctx, sn)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return t, nil
}

// List returns paginated terminals, filtered by store data scope.
func (s *TerminalService) List(ctx context.Context, req dto.TerminalListReq, storeIDs []int64) ([]model.Terminal, int64, error) {
	return s.terminalRepo.List(ctx, req, storeIDs)
}

// Update updates a terminal.
func (s *TerminalService) Update(ctx context.Context, id int64, req dto.UpdateTerminalReq) error {
	t, err := s.terminalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrTerminalNotFound
	}

	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Type != "" {
		t.Type = req.Type
	}
	if req.StoreID != nil {
		store, err := s.storeRepo.GetByID(ctx, *req.StoreID)
		if err != nil {
			return err
		}
		if store == nil {
			return ErrStoreNotFound
		}
		t.StoreID = *req.StoreID
	}

	// Handle status transition
	if req.Status != nil {
		if err := s.changeStatus(ctx, t, *req.Status); err != nil {
			return err
		}
	}

	return s.terminalRepo.Update(ctx, t)
}

// Delete removes a terminal and cleans up Redis.
func (s *TerminalService) Delete(ctx context.Context, id int64) error {
	t, err := s.terminalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrTerminalNotFound
	}

	// Clean up Redis heartbeat key
	s.heartbeatCache.Delete(ctx, t.SN)
	s.heartbeatCache.DeleteDeviceToken(ctx, t.DeviceToken)

	return s.terminalRepo.Delete(ctx, id)
}

// Heartbeat processes a terminal heartbeat.
func (s *TerminalService) Heartbeat(ctx context.Context, sn, ip, mac string) error {
	t, err := s.terminalRepo.GetBySN(ctx, sn)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrTerminalNotFound
	}

	switch t.Status {
	case model.TerminalStatusDisabled:
		return fmt.Errorf("%w: terminal disabled", ErrTerminalDisabled)
	case model.TerminalStatusOffline:
		// Transition offline → online
		s.logRepo.CreateTerminal(ctx, &model.TerminalLog{
			TerminalID: &t.ID,
			SN:         sn,
			EventType:  "online",
		})
	}

	// Update heartbeat
	if err := s.terminalRepo.UpdateHeartbeat(ctx, t.ID, ip); err != nil {
		return err
	}

	// Update MAC if provided
	if mac != "" && t.MACAddress != mac {
		t.MACAddress = mac
		s.terminalRepo.Update(ctx, t)
	}

	// Set Redis heartbeat TTL
	s.heartbeatCache.Set(ctx, sn)

	return nil
}

// RotateToken generates a new device_token for the terminal.
func (s *TerminalService) RotateToken(ctx context.Context, sn string) (string, error) {
	t, err := s.terminalRepo.GetBySN(ctx, sn)
	if err != nil {
		return "", err
	}
	if t == nil {
		return "", ErrTerminalNotFound
	}

	oldToken := t.DeviceToken
	newToken := uuid.New().String()

	if err := s.terminalRepo.UpdateDeviceToken(ctx, t.ID, newToken); err != nil {
		return "", err
	}

	// Invalidate old token in Redis
	s.heartbeatCache.DeleteDeviceToken(ctx, oldToken)

	return newToken, nil
}

// HandleStatusTimeout handles the Redis keyspace notification for heartbeat expiry.
func (s *TerminalService) HandleStatusTimeout(ctx context.Context, sn string) {
	t, err := s.terminalRepo.GetBySN(ctx, sn)
	if err != nil || t == nil {
		return
	}

	if t.Status == model.TerminalStatusOnline {
		s.terminalRepo.UpdateStatusBySN(ctx, sn, model.TerminalStatusOffline)
		s.logRepo.CreateTerminal(ctx, &model.TerminalLog{
			TerminalID: &t.ID,
			SN:         sn,
			EventType:  "heartbeat_timeout",
		})
	}
}

func (s *TerminalService) changeStatus(ctx context.Context, t *model.Terminal, newStatus model.TerminalStatus) error {
	switch newStatus {
	case model.TerminalStatusDisabled:
		if t.Status != model.TerminalStatusDisabled {
			s.heartbeatCache.Delete(ctx, t.SN)
			s.logRepo.CreateTerminal(ctx, &model.TerminalLog{
				TerminalID: &t.ID, SN: t.SN, EventType: "disabled",
			})
		}
		t.Status = model.TerminalStatusDisabled
	case model.TerminalStatusEnabled:
		if t.Status == model.TerminalStatusDisabled {
			t.Status = model.TerminalStatusOffline
			s.logRepo.CreateTerminal(ctx, &model.TerminalLog{
				TerminalID: &t.ID, SN: t.SN, EventType: "enabled",
			})
		}
	default:
		return ErrInvalidStatusChange
	}
	return nil
}

// Sentinel errors for terminal service.
var (
	ErrTerminalNotFound    = errors.New("terminal not found")
	ErrTerminalDisabled    = errors.New("terminal disabled")
	ErrInvalidStatusChange = errors.New("invalid status change")
)
