package documents

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateAccountingDocument(t *testing.T) {
	// given

	// shared between transfers[0] and transfers[1]
	monthLayout := "2006-Jan"
	monthString := "2019-May"
	testMonth, err := time.Parse(monthLayout, monthString)
	if err != nil {
		t.Errorf("Error in test setup - Can't parse testmonth")
	}
	// specific to transfers[0]
	testProjectName := "Doe Company"
	testProjectID := "12345"
	testContactPerson := "John Doe"
	testAmount := 1.44

	// specific to transfers[1]
	testProjectNameAlt := "Doe, Company"
	testProjectIDAlt := "A-B"
	testContactPersonAlt := "Jane Doe"
	testAmountAlt := 1.44

	// for other ChargeBack fields
	testMargin := 0.2
	testCurrency := "USD"
	testProviderResponse := "Manually generated test provider response"
	testTimeStamp := time.Now()

	// expected
	expected := "Position,Month,ProjectName,ProjectID,ContactPerson,Margin,Amount,Currency\n" +
		fmt.Sprintf("1,2019-May,%s,%s,%s,0.200,1.440,%s\n", testProjectName, testProjectID, testContactPerson, testCurrency) +
		fmt.Sprintf("2,2019-May,\"%s\",%s,%s,0.200,1.440,%s\n", testProjectNameAlt, testProjectIDAlt, testContactPersonAlt, testCurrency)

	chargeback := NewChargeBack(
		[]Bill{
			NewBill(
				testProjectName,
				testProjectID,
				testContactPerson,
				testAmount,
				testCurrency,
			),
			NewBill(
				testProjectNameAlt,
				testProjectIDAlt,
				testContactPersonAlt,
				testAmountAlt,
				testCurrency,
			),
		},
		testMargin,
		testMonth,
		testProviderResponse,
		testTimeStamp,
	)

	// when
	accountingDoc := GenerateAccountingDocumentWithLocale(chargeback, "EN")

	if expected != accountingDoc {
		t.Errorf("%v should be %v", accountingDoc, expected)
	}

	// then
}

func TestNewChargeBack(t *testing.T) {
	// given
	monthLayout := "2006-Jan"
	monthString := "2019-May"
	testMonth, err := time.Parse(monthLayout, monthString)
	if err != nil {
		t.Errorf("Error in test setup - Can't parse testmonth")
	}
	testProjectName := "Doe Company"
	testProjectID := "12345"
	testContactPerson := "John Doe"
	testAmount := 1.44
	testCurrency := "USD"
	entry := NewBill(testProjectName, testProjectID, testContactPerson, testAmount, testCurrency)

	entries := []Bill{
		entry,
	}
	testMargin := 0.2

	// when
	chargeback := NewChargeBack(entries, testMargin, testMonth, "test provider response", time.Now())

	// then
	for i, inputEntry := range chargeback.bills {
		if !cmp.Equal(entries[i], inputEntry) {
			t.Fail()
		}
	}
}

func TestNewBill(t *testing.T) {
	testProjectName := "Doe Company"
	testProjectID := "12345"
	testContactPerson := "John Doe"
	testAmount := 1.44
	testCurrency := "USD"
	entry := NewBill(testProjectName, testProjectID, testContactPerson, testAmount, testCurrency)

	if entry.ProjectID != testProjectID ||
		entry.ContactPerson != testContactPerson ||
		entry.Amount != testAmount {
		t.Fail()
	}

}

// TODO: test once more with DecimalCommas on
// Consider: Passing the decimal operator explicitly may be a better solution, as it avoids different control flow branches
func TestFormatRow(t *testing.T) {
	// given
	monthLayout := "2006-Jan"
	monthString := "2019-May"
	testmonth, err := time.Parse(monthLayout, monthString)
	if err != nil {
		t.Errorf("Error in test setup - Can't parse testmonth")
	}

	testpos := 1
	testName := "testName"
	testID := "123"
	testContactPerson := "John Doe"
	testMargin := 0.2
	testAmount := 1.23456789
	testCurrency := "USD"
	expected := []string{"1", monthString, testName, testID, testContactPerson, "0.200", "1.235", testCurrency}

	// when
	row := formatRow(testpos, testmonth, testName, testID, testContactPerson, testMargin, testAmount, testCurrency, "EN")

	// then
	if !cmp.Equal(row, expected) {
		t.Errorf("%v should be %v", row, expected)
	}
}

func TestTrim(t *testing.T) {

	// short text
	s := strings.Repeat("a", 10)
	if trim(s) != s {
		t.Errorf("%s should be %s\n", trim(s), s)
	}

	// empty text
	s = ""
	if trim(s) != s {
		t.Errorf("%s should be %s\n", trim(s), s)
	}

	// too long text
	s = strings.Repeat("a", 51)
	expected := strings.Repeat("a", 47) + "..."
	if trim(s) != expected {
		t.Errorf("%s should be %s\n", trim(s), expected)
	}

}
