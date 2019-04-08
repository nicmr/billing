package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

type CsvEntry struct {
	TimePeriodStart string
	TimePeriodEnd   string
	Amount          string
}

func CreateCsv(csvEntries []CsvEntry) string {
	csvFileContent := new(bytes.Buffer)
	writer := csv.NewWriter(csvFileContent)

	err := writer.Write([]string{"TimePeriodStart", "TimePeriodEnd", "Amount"})
	if err != nil {
		fmt.Print(err)
	}

	for _, entry := range csvEntries {
		value := []string{entry.TimePeriodStart, entry.TimePeriodEnd, entry.Amount}
		err := writer.Write(value)
		if err != nil {
			fmt.Print(err)
		}
	}

	writer.Flush()
	return (csvFileContent.String())
}
