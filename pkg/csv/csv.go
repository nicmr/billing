// This package provides functionality to create .csv-files.

package csv

import (
	"bytes"
	"encoding/csv"
	"log"
)

// Entry is a struct that holds the costs for a specified period of time.
type Entry struct {
	TimePeriodStart string
	TimePeriodEnd   string
	Amount          string
	ID              int
	ProjectID       string
	ContactPerson   string
	Month           string
}

// String parses csvEntries and returns a them as a string with .csv-Formatting,
// i. e. as a csv header followed by the csvEntries.
func String(csvEntries []Entry) string {
	csvFileContent := new(bytes.Buffer)
	writer := csv.NewWriter(csvFileContent)

	err := writer.Write([]string{"TimePeriodStart", "TimePeriodEnd", "Amount", "ID", "ProjectID", "ContactPerson", "Month"})
	if err != nil {
		log.Fatal("Could not write csv header: ", err)
	}

	for _, entry := range csvEntries {
		// Create an []string, the entires will be written comma-serpated to the writer.
		value := []string{entry.TimePeriodStart, entry.TimePeriodEnd, entry.Amount, "someID", "SomeProjectID", "SomePerson", "SomeMonth"}
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Could not write csv entry: ", err)
		}
	}

	writer.Flush()
	return (csvFileContent.String())
}
