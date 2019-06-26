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
