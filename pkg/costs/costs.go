// Package costs offers functions for retrieving cost and usage data
// from the APIs of AWS, Azure and Kubermatic.
package costs

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Altemista/altemista-billing/pkg/csv"
)

// APICallResult contains a Timestamp and ResponseString
// Timestamp is a time.Time of the moment the query was completed.
// ResponseString is a string representation for an arbitrary costexplorer query response
// CsvFileContent is a string with a csv representation of the most important data returned ny the ApiCall
// This struct is likely to change a lot during development, so don't rely too much on its internals.
type APICallResult struct {
	Timestamp      time.Time
	ResponseString string
	// CsvEntries     []csv.Entry
	Entries []Entry
}

// Entry is a struct that holds the relevant returned information returned by an APICall
type Entry struct {
	Month         string `csv:"Month" order:"0"`
	ProjectID     string `csv:"ProjectID" order:"1"`
	ContactPerson string `csv:"ContactPerson" order:"2"`
	Amount        string `csv:"Amount" order:"3"`
	Margin        string `csv:"Margin" order:"4"`
	Total         string `csv:"Total" order:"5"`
}

// ToCsvString returns a csv-conformant string representation of apiResult as a
func (apiResult APICallResult) ToCsvString() string {

	// convert our entry slice to an interface slice
	// (unfortunately go doesn't do this for us)
	slice := make([]interface{}, len(apiResult.Entries))
	for i, v := range apiResult.Entries {
		slice[i] = v
	}

	csvString, err := csv.MarshalReflect(slice)
	if err != nil {
		log.Fatal("unable to generate csv: ", err)
	}
	return csvString
}

// Default returns a default `APICall`, currently AWS
func Default() APICall {
	return AWS()
}

// ApplyMargin applies a margin to the provdided APICallResult
func ApplyMargin(apiResult APICallResult, margin float64) APICallResult {
	for i, entry := range apiResult.Entries {
		amount, err := strconv.ParseFloat(entry.Amount, 64)
		if err != nil {
			log.Fatal("unable to parse cost value in apiResult: ", err)
		}
		total := amount * (1.0 + margin)

		apiResult.Entries[i].Margin = fmt.Sprintf("%v", margin)
		apiResult.Entries[i].Total = fmt.Sprintf("%v", total)
	}
	return apiResult
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
