package parse

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParseCSVFile(t *testing.T) {
	csvFile, err := os.Open("test_data/example.csv")
	if err != nil {
		t.Errorf("error opening example csv file: %s", err)
	}
	transactions, err := ParseCSVFile(csvFile)
	if err != nil {
		t.Errorf("error parsing csv file: %s", err)
	}
	fmt.Println(transactions[0])
	fmt.Println(len(transactions))
}

func TestParseDate(t *testing.T) {
	empty_date, err := ParseDate("")
	if err == nil {
		t.Errorf("expected error when parsing empty date, parsing result: %v", empty_date)
	}
	date_wo_time, err := ParseDate("2022-08-23 ")
	if err == nil {
		t.Errorf("expected error when parsind date w/o time, parsing result: %v", date_wo_time)
	}
	correct_date, err := ParseDate("2022-08-23 9:04:06")
	if err != nil {
		t.Errorf("error parsing correct date: %s", err)
	}
	if correct_date != time.Date(2022, 8, 23, 9, 4, 6, 0, time.UTC) {
		t.Errorf("got wrong date after parsing: %v", correct_date)
	}
}
