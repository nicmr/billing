package costs

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/csv"
)

// costsMonthlyAzure returns the cloud usage costs for the specified month on Azure. Not yet implemented.
func costsMonthlyAzure(month time.Time) (APICallResult, error) {

	result := APICallResult{
		Timestamp:  time.Now(),
		Response:   "not yet implemented for Azure",
		CsvEntries: []csv.Entry{},
	}

	return result, nil
}
