package parse

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/asstronom/EVO_tech_test/pkg/domain"
)

func ParseDate(datestr string) (time.Time, error) {
	var year, month, day, hours, minutes, seconds int
	splits := strings.Split(datestr, " ")
	if len(splits) != 2 {
		return time.Unix(0, 0), fmt.Errorf("wrong date format")
	}
	datesplit := strings.Split(splits[0], "-")
	if len(datesplit) != 3 {
		return time.Unix(0, 0), fmt.Errorf("wrong date format")
	}
	timesplit := strings.Split(splits[1], ":")
	if len(timesplit) != 3 {
		return time.Unix(0, 0), fmt.Errorf("wrong date format")
	}
	year, err := strconv.Atoi(datesplit[0])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing year: %w", err)
	}
	month, err = strconv.Atoi(datesplit[1])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing month: %w", err)
	}
	day, err = strconv.Atoi(datesplit[2])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing day: %w", err)
	}
	hours, err = strconv.Atoi(timesplit[0])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing hours: %w", err)
	}
	minutes, err = strconv.Atoi(timesplit[1])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing minutes: %w", err)
	}
	seconds, err = strconv.Atoi(timesplit[2])
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("error parsing seconds: %w", err)
	}
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.Local), nil
}

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

func ParseCSVFile(file *os.File) ([]domain.Transaction, error) {
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

func ValidatePort(port string) error {
	if len(port) == 0 {
		return fmt.Errorf("port is empty")
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
