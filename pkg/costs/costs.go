// Package query contains the
package costs

import "time"

// APICallResult contains a Timestamp and Response
// Timestamp is a timestamp of the moment the query was completed.
// Response is a string representation for an arbitrary costexplorer query response
// This implementation is a placeholder, when azure and on_premise support has been added,
// outputs from the different queried objects will be parsed into a common result format
type APICallResult struct {
	Timestamp      time.Time
	Response       string
	CsvFileContent string
}

// Default returns a default Type from those implementing CostsQuery
func Default() APICall {
	return AWS()
}

// APICall is a type representing a function that takes two strings representing the start and end of the queried date range
// and returns an APICallResult containing the response from the implemented API
type APICall func(string, string) (APICallResult, error)

// AWS returns an APICall func that will execute an AWS Cost Explorer API call and return APICallResult
func AWS() APICall {
	return costsBetweenAWS
}

// Azure returns an APICall func that will query an Azure Cost Explorer API and return an APICallResult
func Azure() APICall {
	return costsBetweenAzure
}

// OnPremise returns an APICall func that will determine the costs associated with the OnPremise Cloud usage in an APICallResult
func OnPremise() APICall {
	return costsBetweenOnPremise
}
