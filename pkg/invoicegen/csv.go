// This package provides functionality to create .csv-files.
package invoicegen

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// CSV returns a string csv representation of any type that implements ConvertableToGeneratorInput
func CSV(convertable ConvertableToGeneratorInput) string {
	input := convertable.ToInvoiceGenInput()
	return csvGen(input)
}

func csvGen(genInput GeneratorInput) string {
	// Order should match order of row values below in `orderedRowValues`
	orderedColumnHeaders := []string{
		"Month",
		"ProjectID",
		"ContactPerson",
		"Amount",
		"Margin",
		"Total",
	}

	// Init csv writer
	content := new(bytes.Buffer)
	writer := csv.NewWriter(content)

	// write column headers
	writer.Write(orderedColumnHeaders)

	monthFormat := "2006-Jan"

	for _, row := range genInput.Entries {
		// Order should match order of columns above in `orderedColumnHeaders`
		orderedRowValues := []string{
			row.Month.Format(monthFormat),
			row.ProjectID,
			row.ContactPerson,
			fmt.Sprintf("%g", row.Amount),
			fmt.Sprintf("%g", row.Margin),
			fmt.Sprintf("%g", row.Total),
		}
		writer.Write(orderedRowValues)
	}

	writer.Flush()
	return content.String()
}
