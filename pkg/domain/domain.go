package domain

import "time"

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
