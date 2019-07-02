// Package cmd provides the application code for altemista billing subcommands
package cmd

import (
	"log"
	"time"

	"github.com/Altemista/altemista-billing/pkg/billing"
	"github.com/Altemista/altemista-billing/pkg/documents"
	"github.com/Altemista/altemista-billing/pkg/store"
)

// Invoice executes the application code of altemista billing for the invoice subcommand
func Invoice(provider billing.CloudProvider, month time.Time, margin float64, bucket string) error {

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
	upgroup := new(store.UploadGroup)
	errchans := []chan error{
		upgroup.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month),
		upgroup.Upload(accountingDocumentEN, bucket, "invoiceEN", "csv", month),
		upgroup.Upload(auditLog, bucket, "auditLog", "log", month),
	}

	upgroup.Wait()
	for _, channel := range errchans {
		err := <-channel
		if err != nil {
			log.Printf("Failed to upload element: %v\n", err)
		}
	}

	// Print generated accountingDocument to stdout
	log.Println("generated doc:\n" + accountingDocumentEN)

	return nil
}
