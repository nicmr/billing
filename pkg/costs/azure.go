package costs

import "time"

// CostsBetween is a public wrapper around `costsBetween`
// It adds package-level variables as parameters, forwards the function call and adds a timestamp
func costsBetweenAzure(start string, end string) (APICallResult, error) {

	result := APICallResult{
		Timestamp: time.Now(),
		Response:  "not yet implemented for Azure",
	}

	return result, nil
}
