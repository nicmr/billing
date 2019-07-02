// Package documents creates documents from billing data
package documents

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ChargeBack contains all information about a set of chargebacks required by the functions of this package.
type ChargeBack struct {
	month                     time.Time
	bills                     []Bill
	usedMargin                float64
	provider                  string
	providerResponse          string
	providerResponseTimeStamp time.Time
}

// Bill represents a bill for a single project and its associated costs for a single month.
type Bill struct {
	ProjectName   string
	ProjectID     string
	ContactPerson string
	Amount        float64
	Currency      string
}

// NewChargeBack returns a ChargeBack to be passed to the Different methods of the package.
// This is the preferred method of instantiating a struct of this type.
func NewChargeBack(bills []Bill, margin float64, month time.Time, providerResponse string, providerResponseTimeStamp time.Time) ChargeBack {
	return ChargeBack{
		bills:                     bills,
		usedMargin:                margin,
		providerResponse:          providerResponse,
		providerResponseTimeStamp: providerResponseTimeStamp,
		month:                     month,
	}
}

// NewBill returns a new Bill struct with all necessary fields set.
// This is the preferred method of instantiating a struct of this type.
func NewBill(projectname string, projectID string, contactPerson string, amount float64, currency string) Bill {
	return Bill{
		ProjectName:   projectname,
		ProjectID:     projectID,
		ContactPerson: contactPerson,
		Amount:        amount,
		Currency:      currency,
	}

}

// GenerateAuditLog generates a timestamped audit log with the exact response received from the costapi
func GenerateAuditLog(chargeback ChargeBack) string {

	layout := "2006-01-02 Mon 15:04:05 -0700 MST"
	timestamp := chargeback.providerResponseTimeStamp.Format(layout)

	return fmt.Sprintf("Audit log generated on: %v\nCloud provider: %v\n%v", timestamp, chargeback.provider, chargeback.providerResponse)
}

// GenerateAccountingDocumentWithLocale generates a document for accounting.
// Currently, this is a .csv document.
// CSV can to be adapted for different locales to work with different Microsoft Excel locales.
func GenerateAccountingDocumentWithLocale(chargeback ChargeBack, locale string) string {

	// Init csv writer
	content := new(bytes.Buffer)
	writer := csv.NewWriter(content)

	locale = strings.ToUpper(locale)
	switch locale {
	case "DE":
		writer.Comma = ';'
	default:
		// use the default comma separator as in RFC 4180
	}

	// Order should match order of row values below in `orderedRowValues`
	columnHeaders := []string{
		"Position",
		"Month",
		"ProjectName",
		"ProjectID",
		"ContactPerson",
		"Margin",
		"Amount",
		"Currency",
	}

	// write column headers
	writer.Write(columnHeaders)

	for i, bill := range chargeback.bills {
		// Order should match order of columns above in `orderedColumnHeaders`
		orderedRowValues := formatRow(i+1, chargeback.month, bill.ProjectName, bill.ProjectID, bill.ContactPerson, chargeback.usedMargin, bill.Amount, bill.Currency, locale)
		writer.Write(orderedRowValues)
	}

	writer.Flush()

	return content.String()
}

// formatRow formats a single row in the output document and applies locale-specific formatting where necessary
func formatRow(position int, month time.Time, projectName string, projectID string, contactPerson string, usedMargin float64, amount float64, currency string, locale string) []string {

	// format parameters as strings where necessary
	posString := strconv.Itoa(position)
	monthFormat := "2006-Jan"
	monthString := month.Format(monthFormat)
	name := trim(projectName)
	localisedMargin := formatFloatWithLocale(usedMargin, locale)
	localisedAmount := formatFloatWithLocale(amount, locale)

	return []string{
		posString,
		monthString,
		name,
		projectID,
		contactPerson,
		localisedMargin,
		localisedAmount,
		currency,
	}
}

// formatFloatWithLocale formats the passed float for use with the passed locale string.
// currentyl supported locale strings are as follows:
// - "DE" for decimal comma .
// All locales not presest on the list above will apply the default formatting, using a decimal point
func formatFloatWithLocale(value float64, locale string) string {
	locale = strings.ToUpper(locale)

	str := fmt.Sprintf("%.3f", value)

	switch locale {
	case "DE":
		str = strings.Replace(str, ".", ",", -1)
	default:
		// for other locales, the default format using decimal points is acceptable
	}

	return str
}

func trim(text string) string {
	if len(text) > 50 {
		return text[:47] + "..."
	}
	return text
}
