// Package cmd provides the application code for altemista billing subcommands
package cmd

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
	chargeBack, err := billing.CalculateChargeBack(provider, month, margin)
	if err != nil {
		return err
	}

	// Generate documents with required locales and generate audit log
	accountingDocumentEN := documents.GenerateAccountingDocumentWithLocale(chargeBack, "EN")
	accountingDocumentDE := documents.GenerateAccountingDocumentWithLocale(chargeBack, "DE")
	auditLog := documents.GenerateAuditLog(chargeBack)

	// Upload to S3
	filename := "bills/test_costs_"
	fileExtension := ".csv"
	useTimestamp := true
	store.Upload(accountingDocumentEN, bucket, filename, fileExtension, useTimestamp)
	store.Upload(accountingDocumentDE, bucket, filename+"DE", fileExtension, useTimestamp)
	store.Upload(auditLog, bucket, "auditLog", fileExtension, useTimestamp)

	// Print generated accountingDocument to stdout
	log.Println(accountingDocumentEN)

	return nil
}
