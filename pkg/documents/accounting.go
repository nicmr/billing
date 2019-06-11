// Package documents provides functionality to create documents about billing data
package documents

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ChargeBack contains all information about a set of chargebacks required by the functions of this package
type ChargeBack struct {
	transfers  []Transfer
	usedMargin float64
	currency   string
}

// Transfer represents a single entry of entry-specific information in the ChargeBack struct
type Transfer struct {
	ProjectName   string
	ProjectID     string
	ContactPerson string
	Month         time.Time
	Amount        float64
}

// NewChargeBack returns a ChargeBack to be passed to the Different methods of the package
// This is the preferred method of instantiating a struct of this type
func NewChargeBack(transfers []Transfer, margin float64, currency string) ChargeBack {
	return ChargeBack{
		transfers:  transfers,
		usedMargin: margin,
		currency:   currency,
	}
}

// NewTransfer returns a Transfer to be used to create a GeneratorInput struct.
// This is the preferred method of instantiating a struct of this type
func NewTransfer(projectname string, projectID string, month time.Time, contactPerson string, amount float64) Transfer {
	return Transfer{
		ProjectName:   projectname,
		Month:         month,
		ProjectID:     projectID,
		ContactPerson: contactPerson,
		Amount:        amount,
	}

}

// GenerateAccountingDocument creates a string, formatted for use by accountants that will have to enter the invoice data into their company's ERP system
func GenerateAccountingDocument(chargeback ChargeBack) string {

	// Order should match order of row values below in `orderedRowValues`
	columnHeaders := []string{
		"Position",
		"Month",
		"ProjectName",
		"ProjectID",
		"ContactPerson",
		"Margin",
		"Amount",
	}
	// generate csv columns

	// Init csv writer
	content := new(bytes.Buffer)
	writer := csv.NewWriter(content)

	// Init csv writer for csv compatible with Microsoft Excel with German(DE) locale
	contentGER := new(bytes.Buffer)
	writerGER := csv.NewWriter(contentGER)
	writerGER.Comma = ';'

	// write column headers
	writer.Write(columnHeaders)
	writerGER.Write(columnHeaders)

	for i, row := range chargeback.transfers {
		// Order should match order of columns above in `orderedColumnHeaders`
		orderedRowValues := formatRow(i, row.Month, row.ProjectName, row.ProjectID, row.ContactPerson, chargeback.usedMargin, row.Amount, false)
		orderedRowValuesGER := formatRow(i, row.Month, row.ProjectName, row.ProjectID, row.ContactPerson, chargeback.usedMargin, row.Amount, true)
		writer.Write(orderedRowValues)
		writerGER.Write(orderedRowValuesGER)
	}

	writer.Flush()
	writerGER.Flush()
	fmt.Printf("GER csv:\n%v\n", contentGER.String())

	return content.String()
}

func formatRow(position int, month time.Time, projectName string, projectID string, contactPerson string, usedMargin float64, amount float64, decimalComma bool) []string {
	monthFormat := "2006-Jan"

	marginString := fmt.Sprintf("%.3f", usedMargin)
	amountString := fmt.Sprintf("%.3f", amount)

	if decimalComma {
		amountString = strings.Replace(amountString, ".", ",", -1)
		marginString = strings.Replace(marginString, ".", ",", -1)
	}

	return []string{
		strconv.Itoa(position + 1),
		month.Format(monthFormat),
		trim(projectName),
		projectID,
		contactPerson,
		marginString,
		amountString,
	}
}

func trim(text string) string {
	if len(text) > 50 {
		return text[:47] + "..."
	}
	return text
}

// func formatAmount(value float64) string {

// }

// // ConvertableToGeneratorInput allows you to define a function for converting your structs to `GeneratorInput`s
// // This allows them to be passed directly to the functions in the public interface of this package.
// type ConvertableToChargeBack interface {
// 	ToInvoiceGenInput() GeneratorInput
// }
