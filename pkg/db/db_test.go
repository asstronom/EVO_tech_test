package db

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
	"github.com/asstronom/EVO_tech_test/pkg/parse"
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
	TestTransaction1 = domain.Transaction{
		ID:                2,
		RequestID:         20030,
		TerminalID:        3507,
		PartnerObjectID:   1111,
		AmountTotal:       1.0,
		AmountOriginal:    1.0,
		CommisionPs:       0.00,
		CommisionClient:   0.00,
		CommisionProvider: 0.00,
		DateInput:         time.Date(2022, 8, 12, 12, 36, 52, 0, time.Local),
		DatePost:          time.Date(2022, 8, 12, 15, 36, 53, 0, time.Local),
		Status:            "accepted",
		PaymentType:       "cash",
		PaymentNumber:     "PS16698215",
		ServiceID:         13990,
		Service:           "Поповнення карток",
		PayeeID:           14232155,
		PayeeName:         "privat",
		PayeeBankMfo:      255752,
		PayeeBankAccount:  "UA713461333619513",
		PaymentNarrative:  `Перерахування коштів згідно договору про надання послуг А11/27123 від 19.11.2020 р.`,
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
	fmt.Printf("%#v\n", res)
}

func TestInsertTransactionsById(t *testing.T) {
	db, err := Open(context.Background(), dburl)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}
	err = db.InsertTransactions(context.Background(), []domain.Transaction{TestTransaction, TestTransaction1})
	if err != nil {
		t.Errorf("error inserting transactions: %s", err)
	}
}

func TestParseAndInsertCSVFile(t *testing.T) {
	db, err := Open(context.Background(), dburl)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}
	file, err := os.Open("example.csv")
	if err != nil {
		t.Errorf("error opening csv file: %s", err)
	}
	defer file.Close()
	trxs, err := parse.ParseCSVFile(file)
	if err != nil {
		t.Errorf("error parsing csv file: %s", err)
	}
	err = db.InsertTransactions(context.Background(), trxs)
	if err != nil {
		t.Errorf("error inserting many transactions")
	}
}

func TestGetTransactionsWithFilters(t *testing.T) {
	db, err := Open(context.Background(), dburl)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}
	date_from, err := parse.ParseDate("2022-08-17 12:53:44")
	if err != nil {
		t.Errorf("error parsing date_from: %s", err)
	}
	date_to, err := parse.ParseDate("2022-08-17 14:25:27")
	if err != nil {
		t.Errorf("error parsing date_to: %s", err)
	}
	fmt.Println(date_from, date_to)
	filters := map[string]interface{}{
		//"terminal_ids":      []interface{}{3506, 3508, 3511, 3515},
		//"status":       "accepted",
		//"payment_type": "card",
		//"date_post_from":    date_from,
		//"date_post_to":      date_to,
		//"payment_narrative": "27122",
	}
	trxs, err := db.GetTransactions(context.Background(), filters)
	if err != nil {
		t.Errorf("error getting transactions with filters: %s", err)
	}
	fmt.Println("len:", len(trxs))
	for i := range trxs {
		fmt.Println(trxs[i])
	}
}
