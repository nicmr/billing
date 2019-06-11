package billing

// This  file contains functions around the exported types of the billing package:
// CloudProvider and APICALL

import (
	"time"
)

// APICall is a type representing a function that takes a time.Time and returns the AWS costs
// for the month associated with that time.Time
// and returns an APICallResult containing the response from the implemented API
type APICall func(time.Time) (apiCallResult, error)

// CloudProvider wraps the function required to retrieve billing data from a CloudProvider, such as AWS or Azure
type CloudProvider struct {
	apicall APICall
}

// Default returns a default `APICall`, currently AWS
func Default() CloudProvider {
	return AWS()
}

// AWS returns a CloudProvider struct to be passed to CostCalc for cost calculation with the AWS Cost Explorer API
func AWS() CloudProvider {
	return CloudProvider{apicall: costsMonthlyAWS}
}

// Azure returns a CloudProvider struct to be passed to CostCalc for cost calculation with the Azure Cost Explorer API
func Azure() CloudProvider {
	return CloudProvider{apicall: costsMonthlyAzure}
}

// OnPremise returns a CloudProvider struct to be passed to CostCalc for cost calculation with an on-premise solution
func OnPremise() CloudProvider {
	return CloudProvider{apicall: costsMonthlyOnPremise}
}

// apiCallResult contains a Timestamp and ResponseString
// Timestamp is a time.Time of the moment the query was completed.
// ResponseString is a string representation for an arbitrary costexplorer query response
// CsvFileContent is a string with a csv representation of the most important data returned ny the ApiCall
// This struct is likely to change a lot during development, so don't rely too much on its internals.
type apiCallResult struct {
	Timestamp      time.Time
	ResponseString string
	Entries        []apiCallResultEntry
}

// apiCallResultEntry is a struct that holds the relevant returned information of a single entry returned by an APICall
type apiCallResultEntry struct {
	ProjectID string
	Amount    float64
	Currency  string
}
