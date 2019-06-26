package billing

import (
	"time"
)

// costsMonthlyAzure returns the cloud usage costs for the specified month on the on-Premise solution. Not yet implemented.
func costsMonthlyOnPremise(month time.Time) (apiCallResult, error) {

	result := apiCallResult{
		Timestamp:      time.Now(),
		ResponseString: "not yet implemented for OnPremise",
		Entries:        []apiCallResultEntry{},
	}

	return result, nil
}
