package domain

import "time"

const layout = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(layout, csv)
	return err
}

type Transaction struct {
	ID                int
	RequestID         int
	TerminalID        int
	PartnerObjectID   int
	AmountTotal       float64
	AmountOriginal    float64
	CommisionPs       float64
	CommisionClient   float64
	CommisionProvider float64
	DateInput         time.Time
	DatePost          time.Time
	Status            string
	PaymentType       string
	PaymentNumber     string
	ServiceID         int
	Service           string
	PayeeID           int
	PayeeName         string
	PayeeBankMfo      int
	PayeeBankAccount  string
	PaymentNarrative  string
}
