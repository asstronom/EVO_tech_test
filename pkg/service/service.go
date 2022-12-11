package service

import (
	"context"

	"github.com/asstronom/EVO_tech_test/pkg/db"
	"github.com/asstronom/EVO_tech_test/pkg/domain"
)

type Service struct {
	db *db.TransactionDB
}

func NewService(db *db.TransactionDB) (*Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	res, err := s.db.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) GetTransactions(ctx context.Context, filters map[string]interface{}) ([]domain.Transaction, error) {
	trxs, err := s.db.GetTransactions(ctx, filters)
	if err != nil {
		return nil, err
	}
	return trxs, err
}

func (s *Service) InsertTransactions(ctx context.Context, trxs []domain.Transaction) error {
	return s.db.InsertTransactions(ctx, trxs)
}
