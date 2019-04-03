package query

import "time"

// OnPremise is a placeholder struct for managing OnPremise cost calculation logic
type OnPremise struct {
}

// CostsBetween is a public wrapper around `costsBetween`
// It adds package-level variables as parameters, forwards the function call and adds a timestamp
func (OnPremise) CostsBetween(start string, end string) (CostsQueryResult, error) {

	result := CostsQueryResult{
		Timestamp: time.Now(),
		Response:  "not yet implemented for OnPremise",
	}

	return result, nil
}
