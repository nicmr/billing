package query

import "time"

// Azure is a placeholder struct for managing Azure cost explorer calls
type Azure struct {
}

// CostsBetween is a public wrapper around `costsBetween`
// It adds package-level variables as parameters, forwards the function call and adds a timestamp
func (Azure) CostsBetween(start string, end string) (CostsQueryResult, error) {

	result := CostsQueryResult{
		Timestamp: time.Now(),
		Response:  "not yet implemented for Azure",
	}

	return result, nil
}
