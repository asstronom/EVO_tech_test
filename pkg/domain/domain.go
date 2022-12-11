package domain

import (
	"context"
	"io"
	"time"
)

const layout = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(layout, csv)
	return err
}

func (date *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + date.Format(layout) + `"`), nil
}

type Transaction struct {
	ID                int      `csv:"TransactionId"`
	RequestID         int      `csv:"RequestId"`
	TerminalID        int      `csv:"TerminalId"`
	PartnerObjectID   int      `csv:"PartnerObjectId"`
	AmountTotal       float64  `csv:"AmountTotal"`
	AmountOriginal    float64  `csv:"AmountOriginal"`
	CommisionPs       float64  `csv:"CommissionPS"`
	CommisionClient   float64  `csv:"CommissionClient"`
	CommisionProvider float64  `csv:"CommissionProvider"`
	DateInput         DateTime `csv:"DateInput"`
	DatePost          DateTime `csv:"DatePost"`
	Status            string   `csv:"Status"`
	PaymentType       string   `csv:"PaymentType"`
	PaymentNumber     string   `csv:"PaymentNumber"`
	ServiceID         int      `csv:"ServiceId"`
	Service           string   `csv:"Service"`
	PayeeID           int      `csv:"PayeeId"`
	PayeeName         string   `csv:"PayeeName"`
	PayeeBankMfo      int      `csv:"PayeeBankMfo"`
	PayeeBankAccount  string   `csv:"PayeeBankAccount"`
	PaymentNarrative  string   `csv:"PaymentNarrative"`
}

type TransactionService interface {
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	GetTransactions(ctx context.Context, filters map[string]interface{}) ([]Transaction, error)
	InsertTransactions(ctx context.Context, trxs io.Reader) error
}
