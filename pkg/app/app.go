// Package app provides the application code for altemista billing
package app

import (
	"log"
	"time"

	"github.com/Altemista/altemista-billing/pkg/billing"
	"github.com/Altemista/altemista-billing/pkg/documents"
	"github.com/Altemista/altemista-billing/pkg/store"
)

// Run executes the application code of altemista billing
func Run(provider billing.CloudProvider, month time.Time, margin float64, bucket string) error {

	// Call the desired API
	apiresult, err := billing.CalculateCosts(provider, month, margin)
	if err != nil {
		return err
	}

	// Generate the accounting document (csv)
	chargeBack := apiresult.ToChargeBack()
	accountingDocument := documents.GenerateAccountingDocument(chargeBack)

	// Upload to S3
	filename := "bills/test_costs_"
	fileExtension := ".csv"
	useTimestamp := true
	store.Upload(accountingDocument, bucket, filename, fileExtension, useTimestamp)

	// Print generated accountingDocument to stdout
	log.Println(accountingDocument)

	return nil
}
