// This package provides functionality to create .csv-files.

package csv

import (
	"bytes"
	"encoding/csv"
	"log"
)

// Entry is a struct that holds the costs for a specified period of time.
type Entry struct {
	Month         string
	ProjectID     string
	ContactPerson string
	Amount        string
	Margin        string
	Total         string
}

func (e Entry) stringSlice() []string {
	return []string{
		e.Month,
		e.ProjectID,
		e.ContactPerson,
		e.Amount,
		e.Margin,
		e.Total,
	}
}

var csvHeaders = []string{
	"Month",
	"ProjectID",
	"ContactPerson",
	"Amount",
	"Margin",
	"Total",
}

// Marshal parses csvEntries and returns a them as a string with .csv-Formatting,
// i. e. as a csv header followed by the csvEntries.
func Marshal(csvEntries []Entry) string {
	csvFileContent := new(bytes.Buffer)
	writer := csv.NewWriter(csvFileContent)

	err := writer.Write(csvHeaders)
	if err != nil {
		log.Fatal("Could not write csv header: ", err)
	}

	for _, entry := range csvEntries {
		// Create a []string, the entries will be written comma-seperated by the writer.
		row := entry.stringSlice()
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Could not write csv entry: ", err)
		}
	}

	writer.Flush()
	return (csvFileContent.String())
}
