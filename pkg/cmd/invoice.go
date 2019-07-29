// Package cmd provides the application code for altemista billing subcommands
package cmd

import (
	"log"
	"time"

	"github.com/Altemista/altemista-billing/pkg/billing"
	"github.com/Altemista/altemista-billing/pkg/documents"
	"github.com/Altemista/altemista-billing/pkg/store"
)

// Invoice executes the application code of altemista billing for the invoice subcommand.
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

	upgroup.Upload(accountingDocumentDE, bucket, "invoiceDE", month)
	upgroup.Upload(accountingDocumentEN, bucket, "invoiceEN", month)
	upgroup.Upload(auditLog, bucket, "auditLog", month)

	upgroup.Wait()

	for _, out := range upgroup.Outputs {
		err := <-out.Err
		if err != nil {
			log.Printf("Failed to upload element: %v\n", err)
		} else {
			s3output := <-out.S3Output
			if s3output == nil {
				log.Println("Output of S3 was nil, can't display upload information")
			} else {
				log.Printf("Uploaded as %v\n", s3output.Location)
			}
		}
	}

	// Print generated accountingDocument to stdout
	log.Println("generated Document:\n" + accountingDocumentEN)

	return nil
}
