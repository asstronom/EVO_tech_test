package parse

import (
	"fmt"
	"os"
	"strconv"
	"testing"
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
	fmt.Print(transactions[0])
}

func TestParseDate(t *testing.T) {
	time, err := ParseDate("2022-08-23 9:04:06")
	if err != nil {
		t.Errorf("error parsing time: %s", err)
	}
	fmt.Println(time)
}

func TestParseFloat(t *testing.T) {
	f, err := strconv.ParseFloat("1.00", 64)
	if err != nil {
		t.Errorf("error parsing float: %s", err)
	}
	fmt.Println(f)
}
