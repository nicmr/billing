package costs

import "time"

// costsMonthlyAzure returns the cloud usage costs for the specified month on the on-Premise solution. Not yet implemented.
func costsMonthlyOnPremise(month time.Time) (APICallResult, error) {

	result := APICallResult{
		Timestamp:      time.Now(),
		Response:       "not yet implemented for OnPremise",
		CsvFileContent: "not yet implemented for OnPremise",
	}

	return result, nil
}
