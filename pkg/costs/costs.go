// Package costs offers functions for retrieving cost and usage data
// from the APIs of AWS, Azure and Kubermatic.
package costs

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/invoicegen"
)

// CostCalcResult is a struct that holds all information exposed via calls to the public interface of CostCalc
type CostCalcResult struct {
	Month         time.Time
	Margin        float64
	apicallresult apiCallResult
	Totals        []float64 // index of total corresponds to index of apiCallResult.Entries(?)
}

// CostCalc calculates the costs in the passed month at the specified provider and calculates the total cost with the specified margin.
func CostCalc(provider Provider, month time.Time, margin float64) (CostCalcResult, error) {
	apiresult, err := provider.apicall(month)
	if err != nil {
		return CostCalcResult{}, err
	}

	costcalcresult := CostCalcResult{
		Month:         month,
		Margin:        margin,
		apicallresult: apiresult,
		Totals:        applyMargin(apiresult.Entries, margin),
	}

	return costcalcresult, nil
}

// ToInvoiceGenInput creates a `invoicegen.GeneratorInput` from a `CostCalcresult`
func (costResult CostCalcResult) ToInvoiceGenInput() invoicegen.GeneratorInput {
	length := len(costResult.apicallresult.Entries)
	entries := make([]invoicegen.GeneratorEntry, length)

	for i, entry := range costResult.apicallresult.Entries {
		entries[i] = invoicegen.GeneratorEntry{
			Month:         costResult.Month,
			ProjectID:     entry.ProjectID,
			ContactPerson: "not yet implemented",
			Amount:        entry.Amount,
			Margin:        costResult.Margin,
		}
	}

	return invoicegen.GeneratorInput{Entries: entries}
}

// ApplyMargin applies a margin to the provdided APICallResult
func applyMargin(entries []apiCallResultEntry, margin float64) []float64 {
	results := make([]float64, len(entries))

	for i, entry := range entries {
		total := entry.Amount * (1.0 + margin)
		results[i] = total
	}
	return results
}

// APICall is a type representing a function that takes a time.Time and returns the AWS costs
// for the month associated with that time.Time
// and returns an APICallResult containing the response from the implemented API
type APICall func(time.Time) (apiCallResult, error)

// Provider wraps the function required to retrieve billing data from a provider, such as AWS or Azure
type Provider struct {
	apicall APICall
}

// Default returns a default `APICall`, currently AWS
func Default() Provider {
	return AWS()
}

// AWS returns a Provider struct to be passed to CostCalc for cost calculation with the AWS Cost Explorer API
func AWS() Provider {
	return Provider{apicall: costsMonthlyAWS}
}

// Azure returns a Provider struct to be passed to CostCalc for cost calculation with the Azure Cost Explorer API
func Azure() Provider {
	return Provider{apicall: costsMonthlyAzure}
}

// OnPremise returns a Provider struct to be passed to CostCalc for cost calculation with an on-premise solution
func OnPremise() Provider {
	return Provider{apicall: costsMonthlyOnPremise}
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
