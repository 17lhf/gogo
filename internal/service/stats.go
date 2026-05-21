package service

import (
	"context"

	"gogo/internal/dto"
	"gogo/internal/repository"
)

type StatsService struct {
	terminalRepo repository.TerminalRepository
}

func NewStatsService(terminalRepo repository.TerminalRepository) *StatsService {
	return &StatsService{terminalRepo: terminalRepo}
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
