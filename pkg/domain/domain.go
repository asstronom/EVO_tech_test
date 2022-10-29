package domain

import "time"

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
	PayeeBankAccount  int
	PaymentNarrative  string
}
