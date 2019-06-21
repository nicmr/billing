// Package billing offers functions for retrieving cost and usage data
// from the APIs of AWS, Azure and Kubermatic.
package billing

import (
	"time"

	"github.com/Altemista/altemista-billing/pkg/documents"
)

// CalculateChargeBack calculates the costs in the passed month at the specified provider and calculates the total cost with the specified margin.
func CalculateChargeBack(provider CloudProvider, month time.Time, margin float64) (documents.ChargeBack, error) {

	apiresult, err := provider.apicall(month)
	if err != nil {
		return documents.ChargeBack{}, err
	}

	bills := make([]documents.Bill, len(apiresult.Entries))
	for i, entry := range apiresult.Entries {
		amount := applyMargin(entry, margin)
		bills[i] = documents.NewBill("Project Name not yet implemented", entry.ProjectID, "Contact person not yet implemented", amount, entry.Currency)
	}
	chargeback := documents.NewChargeBack(bills, margin, month, apiresult.ResponseString, apiresult.Timestamp)
	return chargeback, nil
}

// ApplyMargin applies a margin to the provdided APICallResult
func applyMargin(entry apiCallResultEntry, margin float64) float64 {
	total := entry.Amount * (1.0 + margin)
	return total
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

// // ToChargeBack allows conversion to to a ChargeBack struct, and can thus be parsed to functions exported by documents
// func (costResult CostCalcResult) ToChargeBack() documents.ChargeBack {
// 	length := len(costResult.apicallresult.Entries)
// 	entries := make([]documents.Transfer, length)

// 	for i, entry := range costResult.apicallresult.Entries {
// 		entries[i] = documents.NewTransfer(
// 			"name not yet implemented",
// 			entry.ProjectID,
// 			costResult.Month,
// 			"contactperson not yet implemented",
// 			costResult.Totals[i],
// 		)
// 	}

// 	return documents.NewChargeBack(entries, costResult.Margin, costResult.Currency)
// }
