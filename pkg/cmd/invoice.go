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

	// synchronous
	// _, err = store.Upload(accountingDocumentEN, bucket, "invoiceEN", "csv", month)
	// if err != nil {
	// 	log.Println("failed to upload invoiceEN")
	// }
	// _, err = store.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month)
	// if err != nil {
	// 	log.Println("failed to upload invoiceDE")
	// }
	// _, err = store.Upload(auditLog, bucket, "auditLog", "log", month)
	// if err != nil {
	// 	log.Println("failed to upload auditLog")
	// }

	// // asynchronous, but verbose and error-prone (wrong param to wg.Add possible)
	// wg := new(sync.WaitGroup)
	// wg.Add(3)
	// go func() {
	// 	defer wg.Done()
	// 	_, err = store.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month)
	// 	if err != nil {
	// 		log.Println("failed to upload invoiceDE")
	// 	}
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	_, err = store.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month)
	// 	if err != nil {
	// 		log.Println("failed to upload invoiceDE")
	// 	}
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	_, err = store.Upload(auditLog, bucket, "auditLog", "log", month)
	// 	if err != nil {
	// 		log.Println("failed to upload invoiceDE")
	// 	}
	// }()

	// asynchronous, nice interface
	ugroup := new(store.UploadGroup)
	errchans := []chan error{
		ugroup.Upload(accountingDocumentDE, bucket, "invoiceDE", "csv", month),
		ugroup.Upload(accountingDocumentEN, bucket, "invoiceEN", "csv", month),
		ugroup.Upload(auditLog, bucket, "auditLog", "log", month),
	}

	ugroup.Wait()

	for _, ec := range errchans {
		err := <-ec
		if err != nil {
			log.Printf("Failed to upload element: %v\n", err)
		}
	}

	// Print generated accountingDocument to stdout
	log.Println("generated doc:\n" + accountingDocumentEN)

	return nil
}
