package parse

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
)

const layout = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(layout, csv)
	return err
}

type transactionParse struct {
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

// parses date string, format is 'YYYY-MM-DD HH-MM-SS'
func ParseDate(datestr string) (time.Time, error) {
	t, err := time.Parse(layout, datestr)
	if err != nil {
		return time.UnixMilli(0), err
	}
	return t, nil
}


// parses csv file
func ParseCSVFile(file io.Reader) ([]domain.Transaction, error) {
	transactions := []domain.Transaction{}
	err := gocsv.Unmarshal(file, &transactions)
	if err != nil {
		return nil, fmt.Errorf("error parsing csv file: %w", err)
	}
	return transactions, nil
}

// valides port string
func ValidatePort(port string) error {
	if len(port) == 0 {
		return errors.New("port is wrong")
	}
	if port[0] != ':' {
		return fmt.Errorf(`wrong port, first character is not ":"`)
	}
	if len(port) < 2 {
		return fmt.Errorf("port number is not specified")
	}
	n, err := strconv.Atoi(port[1:])
	if err != nil {
		return fmt.Errorf("error port number: %w", err)
	}
	if n < 1 || n > 65535 {
		return fmt.Errorf("port number is not in range 1-65535")
	}
	return nil
}
