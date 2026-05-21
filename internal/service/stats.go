package service

import (
	"context"

	"gogo/internal/dto"
	"gogo/internal/model"
	"gogo/internal/repository"
)

type StatsService struct {
	terminalRepo repository.TerminalRepository
	userRepo     repository.UserRepository
}

func NewStatsService(terminalRepo repository.TerminalRepository, userRepo repository.UserRepository) *StatsService {
	return &StatsService{terminalRepo: terminalRepo, userRepo: userRepo}
}

func (s *StatsService) GetTerminals(ctx context.Context) (*dto.StatsTerminalsResp, error) {
	statsByStatus, err := s.terminalRepo.GetCountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	statsByStore, err := s.terminalRepo.GetCountByStore(ctx)
	if err != nil {
		return nil, err
	}

	statsByRecentAdded, err := s.terminalRepo.GetCountByRecentAdded(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.StatsTerminalsResp{
		StatusDistribution: *statsByStatus,
		ByStore:            statsByStore,
		RecentAdded:        *statsByRecentAdded,
	}, nil
}

func (s *StatsService) GetUsers(ctx context.Context) (*dto.UserStatsResp, error) {
	statusMap, err := s.userRepo.GetCountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	statusDistribution := make(map[string]int64)
	for status, count := range statusMap {
		statusDistribution[model.UserStatus(status).String()] = count
	}

	byRole, err := s.userRepo.GetCountByRole(ctx)
	if err != nil {
		return nil, err
	}

	recentAdded, err := s.userRepo.GetCountByRecentAdded(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.UserStatsResp{
		StatusDistribution: statusDistribution,
		ByRole:             byRole,
		RecentAdded:        *recentAdded,
	}, nil
}
