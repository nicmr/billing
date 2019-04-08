// Package query contains the
package costs

import "time"

// CostsQuery is an interface that contains the `CostsBetween` method.
// CostsBetween allows a type to be queried for the cost between two dates
type CostsQuery interface {
	CostsBetween(string, string) (CostsQueryResult, error)
}

// CostsQueryResult contains a Timestamp and Response
// Timestamp is a timestamp of the moment the query was completed.
// Response is a string representation for an arbitrary costexplorer query response
// This implementation is a placeholder, when azure and on_premise support has been added,
// outputs from the different queried objects will be parsed into a common result format
type CostsQueryResult struct {
	Timestamp      time.Time
	Response       string
	CsvFileContent string
}

// DefaultClient returns a default Type from those implementing CostsQuery
func DefaultClient() CostsQuery {
	return AWS{}
}
