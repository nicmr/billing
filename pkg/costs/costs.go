// Package costs offers functions for retrieving cost and usage data
// from the APIs of AWS, Azure and Kubermatic.
package costs

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/csv"
)

// APICallResult contains a Timestamp and Response
// Timestamp is a time.Time of the moment the query was completed.
// Response is a string representation for an arbitrary costexplorer query response
// CsvFileContent is a string with a csv representation of the most important data returned ny the ApiCall
// This struct is likely to change a lot during development, so don't rely too much on its internals.
type APICallResult struct {
	Timestamp  time.Time
	Response   string
	CsvEntries []csv.Entry
}

// Default returns a default `APICall`, currently AWS
func Default() APICall {
	return AWS()
}

// APICall is a type representing a function that takes a time.Time and returns the AWS costs
// for the month associated with that time.Time
// and returns an APICallResult containing the response from the implemented API
type APICall func(time.Time) (APICallResult, error)

// AWS returns an APICall func that will execute an AWS Cost Explorer API call and return APICallResult
func AWS() APICall {
	return costsMonthlyAWS
}

// Azure returns an APICall func that will query an Azure Cost Explorer API and return an APICallResult
func Azure() APICall {
	return costsMonthlyAzure
}

// OnPremise returns an APICall func that will determine the costs associated with the OnPremise Cloud usage in an APICallResult
func OnPremise() APICall {
	return costsMonthlyOnPremise
}
