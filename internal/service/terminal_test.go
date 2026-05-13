package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gogo/internal/dto"
	"gogo/internal/model"
)

// terminalRepoStub is a stub for TerminalRepository.
type terminalRepoStub struct {
	terminals map[int64]*model.Terminal
	bySN      map[string]*model.Terminal
	nextID    int64
}

func newTerminalRepoStub() *terminalRepoStub {
	return &terminalRepoStub{
		terminals: make(map[int64]*model.Terminal),
		bySN:      make(map[string]*model.Terminal),
		nextID:    1,
	}
}

func (s *terminalRepoStub) Create(ctx context.Context, t *model.Terminal) error {
	t.ID = s.nextID
	s.nextID++
	s.terminals[t.ID] = t
	s.bySN[t.SN] = t
	return nil
}

func (s *terminalRepoStub) GetByID(ctx context.Context, id int64) (*model.Terminal, error) {
	t, ok := s.terminals[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (s *terminalRepoStub) GetBySN(ctx context.Context, sn string) (*model.Terminal, error) {
	t, ok := s.bySN[sn]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (s *terminalRepoStub) List(ctx context.Context, req dto.TerminalListReq, storeIDs []int64) ([]model.Terminal, int64, error) {
	var result []model.Terminal
	for _, t := range s.terminals {
		result = append(result, *t)
	}
	return result, int64(len(result)), nil
}

func (s *terminalRepoStub) Update(ctx context.Context, t *model.Terminal) error {
	s.terminals[t.ID] = t
	s.bySN[t.SN] = t
	return nil
}

func (s *terminalRepoStub) Delete(ctx context.Context, id int64) error {
	t := s.terminals[id]
	delete(s.terminals, id)
	if t != nil {
		delete(s.bySN, t.SN)
	}
	return nil
}

func (s *terminalRepoStub) UpdateHeartbeat(ctx context.Context, id int64, ip string) error {
	t := s.terminals[id]
	if t != nil {
		t.Status = "online"
		t.IPAddress = ip
	}
	return nil
}

func (s *terminalRepoStub) UpdateStatus(ctx context.Context, id int64, status string) error {
	t := s.terminals[id]
	if t != nil {
		t.Status = status
	}
	return nil
}

func (s *terminalRepoStub) UpdateStatusBySN(ctx context.Context, sn string, status string) error {
	t := s.bySN[sn]
	if t != nil {
		t.Status = status
	}
	return nil
}

func (s *terminalRepoStub) UpdateDeviceToken(ctx context.Context, id int64, token string) error {
	t := s.terminals[id]
	if t != nil {
		t.DeviceToken = token
	}
	return nil
}

// storeRepoStub is a minimal stub for StoreRepository.
type storeRepoStub struct {
	stores map[int64]*model.Store
}

func newStoreRepoStub() *storeRepoStub {
	return &storeRepoStub{stores: make(map[int64]*model.Store)}
}

func (s *storeRepoStub) Create(ctx context.Context, store *model.Store) error { return nil }
func (s *storeRepoStub) GetByID(ctx context.Context, id int64) (*model.Store, error) {
	store, ok := s.stores[id]
	if !ok {
		return nil, nil
	}
	return store, nil
}
func (s *storeRepoStub) List(ctx context.Context, req dto.StoreListReq) ([]model.Store, int64, error) {
	return nil, 0, nil
}
func (s *storeRepoStub) Update(ctx context.Context, store *model.Store) error   { return nil }
func (s *storeRepoStub) Delete(ctx context.Context, id int64) error              { return nil }
func (s *storeRepoStub) HasTerminals(ctx context.Context, storeID int64) (bool, error) {
	return false, nil
}

// logRepoStub is a minimal stub for LogRepository.
type logRepoStub struct{}

func newLogRepoStub() *logRepoStub {
	return &logRepoStub{}
}

func (s *logRepoStub) CreateOperation(ctx context.Context, log *model.OperationLog) error { return nil }
func (s *logRepoStub) ListOperations(ctx context.Context, req dto.OperationLogListReq) ([]model.OperationLog, int64, error) {
	return nil, 0, nil
}
func (s *logRepoStub) DeleteOperationBefore(ctx context.Context, before time.Time) error { return nil }
func (s *logRepoStub) CreateTerminal(ctx context.Context, log *model.TerminalLog) error     { return nil }
func (s *logRepoStub) ListTerminals(ctx context.Context, req dto.TerminalLogListReq) ([]model.TerminalLog, int64, error) {
	return nil, 0, nil
}
func (s *logRepoStub) DeleteTerminalBefore(ctx context.Context, before time.Time) error { return nil }

func TestTerminalService_Create(t *testing.T) {
	storeRepo := newStoreRepoStub()
	storeRepo.stores[1] = &model.Store{ID: 1, Name: "Test Store", Code: "SH001"}

	svc := NewTerminalService(newTerminalRepoStub(), storeRepo, nil, newLogRepoStub())

	tm, err := svc.Create(context.Background(), dto.CreateTerminalReq{
		SN:      "TM001",
		Name:    "Test Terminal",
		Type:    "POS",
		StoreID: 1,
	})
	require.NoError(t, err)
	assert.Equal(t, "TM001", tm.SN)
	assert.Equal(t, "offline", tm.Status)
	assert.NotEmpty(t, tm.DeviceToken)
}

func TestTerminalService_Create_StoreNotFound(t *testing.T) {
	svc := NewTerminalService(newTerminalRepoStub(), newStoreRepoStub(), nil, newLogRepoStub())

	_, err := svc.Create(context.Background(), dto.CreateTerminalReq{
		SN:      "TM001",
		Name:    "Test",
		StoreID: 999,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "门店不存在")
}

func TestTerminalService_RotateToken(t *testing.T) {
	repo := newTerminalRepoStub()
	repo.bySN["TM001"] = &model.Terminal{ID: 1, SN: "TM001", DeviceToken: "old-token", Status: "online"}
	repo.terminals[1] = repo.bySN["TM001"]

	// RotateToken requires Redis heartbeatCache - skip for now (integration test needed)
	t.Skip("requires Redis heartbeat cache")
}

func TestTerminalService_Delete(t *testing.T) {
	// Delete requires Redis heartbeatCache - skip for now (integration test needed)
	t.Skip("requires Redis heartbeat cache")
}
