package costs

import (
	"time"
)

// costsMonthlyAzure returns the cloud usage costs for the specified month on Azure. Not yet implemented.
func costsMonthlyAzure(month time.Time) (APICallResult, error) {

	result := APICallResult{
		Timestamp:      time.Now(),
		ResponseString: "not yet implemented for Azure",
		Entries:        []Entry{},
	}

	return result, nil
}
