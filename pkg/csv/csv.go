package csv

import (
	"os"
	"encoding/csv"
	"fmt"
)

type CsvEntry struct {
	TimePeriodStart string
	TimePeriodEnd string
	Amount string
}

func CreateCsv (csvEntries []CsvEntry) {
	// fmt.Print(csvEntries)

	file, err := os.Create("billing.csv")
    if err != nil {
		fmt.Print(err)
	}
    defer file.Close()

    writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"TimePeriodStart", "TimePeriodEnd", "Amount"})
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
}
