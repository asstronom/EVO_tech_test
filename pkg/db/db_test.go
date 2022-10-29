package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
)

const (
	dburl = "postgres://user:mypassword@localhost:5432/transactions"
)

var (
	TestTransaction = domain.Transaction{
		ID:                1,
		RequestID:         20020,
		TerminalID:        3506,
		PartnerObjectID:   1111,
		AmountTotal:       1.0,
		AmountOriginal:    1.0,
		CommisionPs:       0.00,
		CommisionClient:   0.00,
		CommisionProvider: 0.00,
		DateInput:         time.Date(2022, 8, 12, 11, 25, 27, 0, time.Local),
		DatePost:          time.Date(2022, 8, 12, 14, 25, 27, 0, time.Local),
		Status:            "accepted",
		PaymentType:       "cash",
		PaymentNumber:     "PS16698205",
		ServiceID:         13980,
		Service:           "Поповнення карток",
		PayeeID:           14232155,
		PayeeName:         "pumb",
		PayeeBankMfo:      254751,
		PayeeBankAccount:  "UA713451373919523",
		PaymentNarrative:  `Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р.`,
	}
)

func TestInsertSingleTransaction(t *testing.T) {
	db, err := Open(context.Background(), dburl)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}
	err = db.InsertTransaction(context.Background(), TestTransaction)
	if err != nil {
		t.Errorf("error while inserting: %s", err)
	}
}

func TestGetTransactionById(t *testing.T) {
	db, err := Open(context.Background(), dburl)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}
	res, err := db.GetTransactionByID(context.Background(), 1)
	if err != nil {
		t.Errorf("error getting trx: %s", err)
	}
	fmt.Print(res)
}
