package parse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

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

// parses date string, format is 'YYYY-MM-DD HH-MM-SS'
func ParseDate(datestr string) (time.Time, error) {
	t, err := time.Parse(layout, datestr)
	if err != nil {
		return time.UnixMilli(0), err
	}
	return t, nil
}

// converts csv record to domain.Transaction
func recordToTransaction(record []string) (*domain.Transaction, error) {
	var err error
	trx := domain.Transaction{}
	trx.ID, err = strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.RequestID, err = strconv.Atoi(record[1])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.TerminalID, err = strconv.Atoi(record[2])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PartnerObjectID, err = strconv.Atoi(record[3])
	if err != nil {
		return nil, fmt.Errorf("error parsing PartnerObjectID: %w", err)
	}
	trx.AmountTotal, err = strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing AmountTotal: %w", err)
	}
	trx.AmountOriginal, err = strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing AmountOriginal: %w", err)
	}
	trx.CommisionPs, err = strconv.ParseFloat(record[6], 64)
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.CommisionClient, err = strconv.ParseFloat(record[7], 64)
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.CommisionProvider, err = strconv.ParseFloat(record[8], 64)
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.DateInput, err = ParseDate(record[9])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.DatePost, err = ParseDate(record[10])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.Status = record[11]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PaymentType = record[12]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PaymentNumber = record[13]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.ServiceID, err = strconv.Atoi(record[14])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.Service = record[15]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PayeeID, err = strconv.Atoi(record[16])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PayeeName = record[17]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PayeeBankMfo, err = strconv.Atoi(record[18])
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PayeeBankAccount = record[19]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	trx.PaymentNarrative = record[20]
	if err != nil {
		return nil, fmt.Errorf("error converting record to trx: %w", err)
	}
	return &trx, nil
}

// parses csv file
func ParseCSVFile(file io.Reader) ([]domain.Transaction, error) {
	r := csv.NewReader(file)
	columns, err := r.Read()
	fmt.Println(columns)
	if err != nil {
		return nil, fmt.Errorf("error parsing columns of csv file: %w", err)
	}
	res := []domain.Transaction{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error parsing record of csv file: %w", err)
		}
		tx, err := recordToTransaction(record)
		if err != nil {
			return nil, fmt.Errorf("error converting record to Transaction: %w", err)
		}
		res = append(res, *tx)
	}
	return res, nil
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
