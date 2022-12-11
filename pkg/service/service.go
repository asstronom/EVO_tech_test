package service

import (
	"context"
	"fmt"
	"io"

	"github.com/asstronom/EVO_tech_test/pkg/db"
	"github.com/asstronom/EVO_tech_test/pkg/domain"
	"github.com/gocarina/gocsv"
)

const (
	batchSize = 40
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

func (s *Service) InsertTransactions(ctx context.Context, trxs io.Reader) error {
	//recieve values one by one using channel
	trxchan := make(chan domain.Transaction)
	go gocsv.UnmarshalToChan(trxs, trxchan)

	//send transactions in batches
	transactions := make([]domain.Transaction, 0, batchSize)
	for trx := range trxchan {
		transactions = append(transactions, trx)
		//if batch is full - send it
		if len(transactions) == batchSize {
			err := s.db.InsertTransactions(ctx, transactions)
			if err != nil {
				return fmt.Errorf("error inserting transactions: %w", err)
			}
			transactions = transactions[:0]
		}
	}
	//send leftovers
	if len(transactions) != 0 {
		err := s.db.InsertTransactions(ctx, transactions)
		if err != nil {
			return fmt.Errorf("error inserting leftover transactions: %w", err)
		}
	}
	return nil
}
