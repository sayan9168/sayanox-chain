package storage

import (
	"context"

	"github.com/sayan9168/sayanox-chain/internal/config"
)

type Service struct {
	cfg *config.Config
	db  *LevelDB
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) Name() string {
	return "Storage"
}

func (s *Service) Start(ctx context.Context) error {
	db, err := Open(s.cfg.DataDir)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Service) DB() *LevelDB {
	return s.db
}
