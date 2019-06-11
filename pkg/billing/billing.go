// Package billing offers functions for retrieving cost and usage data
// from the APIs of AWS, Azure and Kubermatic.
package billing

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/documents"
)

// CostCalcResult is a struct that holds all information exposed via calls to the public interface of CostCalc
type CostCalcResult struct {
	Month         time.Time
	Margin        float64
	apicallresult apiCallResult
	Totals        []float64 // index of total corresponds to index of apiCallResult.Entries(?)
	Currency      string
}

// CalculateCosts calculates the costs in the passed month at the specified provider and calculates the total cost with the specified margin.
func CalculateCosts(provider CloudProvider, month time.Time, margin float64) (CostCalcResult, error) {
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

// // ToInvoiceGenInput creates a `invoicegen.GeneratorInput` from a `CostCalcresult`
// func (costResult CostCalcResult) ToInvoiceGenInput() invoicegen.GeneratorInput {
// 	length := len(costResult.apicallresult.Entries)
// 	entries := make([]invoicegen.GeneratorEntry, length)

// 	for i, entry := range costResult.apicallresult.Entries {
// 		projectid := entry.ProjectID
// 		if projectid == "" {
// 			projectid = "no Project ID assigned"
// 		}
// 		entries[i] = invoicegen.NewGeneratorEntry(
// 			costResult.Month,
// 			projectid,
// 			"not yet implemented",
// 			entry.Amount,
// 			costResult.Margin,
// 			costResult.Totals[i],
// 		)
// 	}

// 	return invoicegen.NewGeneratorInput(entries)
// }

// ToChargeBack allows conversion to to a ChargeBack struct, and can thus be parsed to functions exported by documents
func (costResult CostCalcResult) ToChargeBack() documents.ChargeBack {
	length := len(costResult.apicallresult.Entries)
	entries := make([]documents.Transfer, length)

	for i, entry := range costResult.apicallresult.Entries {
		entries[i] = documents.NewTransfer(
			"name not yet implemented",
			entry.ProjectID,
			costResult.Month,
			"contactperson not yet implemented",
			costResult.Totals[i],
		)
	}

	return documents.NewChargeBack(entries, costResult.Margin, costResult.Currency)
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
