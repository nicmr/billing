package billing

// This file contains functions around the exported types of the billing package:
// CloudProvider and APICall

import (
	"time"
)

// APICall is a type representing a function that takes a time.Time and returns the and returns an APICallResult containing
// the costs for the month associated with that time.Time.
type APICall func(time.Time) (apiCallResult, error)

// CloudProvider wraps the function required to retrieve billing data from a CloudProvider, such as AWS or Azure
type CloudProvider struct {
	apicall APICall
	name    string
}

// Default returns a default `APICall`, currently AWS
func Default() CloudProvider {
	return AWS()
}

// AWS returns a CloudProvider struct that implements an APICall  the AWS Cost Explorer API
func AWS() CloudProvider {
	return CloudProvider{apicall: costsMonthlyAWS, name: "AWS"}
}

// Azure returns a CloudProvider struct  the Azure Cost Explorer API
func Azure() CloudProvider {
	return CloudProvider{apicall: costsMonthlyAzure, name: "Azure"}
}

// OnPremise returns a CloudProvider struct implements an APICall for an Altemista on-premise solution
func OnPremise() CloudProvider {
	return CloudProvider{apicall: costsMonthlyOnPremise, name: "OnPremise"}
}

// apiCallResult contains a Timestamp and ResponseString.
// Timestamp is a time.Time of the moment the query was completed.
// the Response string is the exact response received from the apiCall.
// Entries is a parsed Representation of the individual entries of the reply.
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
