// This package provides functionality to create .csv-files.

package csv

import (
	"bytes"
	"encoding/csv"
	"log"
)

// CsvEntry is a struct that holds the costs for a specified period of time.
type CsvEntry struct {
	TimePeriodStart string
	TimePeriodEnd   string
	Amount          string
}

// CreateCsv is a function that creates (the content of) a .csv-file containing
// a header followed by the csvEntries. The content is returned as a string. This
// function does NOT create an actual file in the file system.
func CreateCsv(csvEntries []CsvEntry) string {
	csvFileContent := new(bytes.Buffer)
	writer := csv.NewWriter(csvFileContent)

	err := writer.Write([]string{"TimePeriodStart", "TimePeriodEnd", "Amount"})
	if err != nil {
		log.Fatal("Could not write csv header: ", err)
	}

	for _, entry := range csvEntries {
		// Create an []string, the entires will be written comma-serpated to the writer.
		value := []string{entry.TimePeriodStart, entry.TimePeriodEnd, entry.Amount}
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Could not write csv entry: ", err)
		}
	}

	writer.Flush()
	return (csvFileContent.String())
}
