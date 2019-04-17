package costs

import "time"

// CostsBetween is a public wrapper around `costsBetween`
// It adds package-level variables as parameters, forwards the function call and adds a timestamp
func costsMonthlyOnPremise(month time.Time) (APICallResult, error) {

	result := APICallResult{
		Timestamp:      time.Now(),
		Response:       "not yet implemented for OnPremise",
		CsvFileContent: "not yet implemented for OnPremise",
	}

	return result, nil
}
