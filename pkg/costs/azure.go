package costs

import (
	"time"
)

// costsMonthlyAzure returns the cloud usage costs for the specified month on Azure. Not yet implemented.
func costsMonthlyAzure(month time.Time) (apiCallResult, error) {

	result := apiCallResult{
		Timestamp:      time.Now(),
		ResponseString: "not yet implemented for Azure",
		Entries:        []apiCallResultEntry{},
	}

	return result, nil
}
