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
	store.Upload(accountingDocumentEN, bucket, "invoice", "csv", month)
	store.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month)
	store.Upload(auditLog, bucket, "auditLog", "log", month)

	// Print generated accountingDocument to stdout
	log.Println(accountingDocumentEN)

	return nil
}
