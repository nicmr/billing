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
	testCurrency := "$"

	// expected
	expected := "Position,Month,ProjectName,ProjectID,ContactPerson,Margin,Amount\n" +
		fmt.Sprintf("1,2019-May,%s,%s,%s,0.200,1.440\n", testProjectName, testProjectID, testContactPerson) +
		fmt.Sprintf("2,2019-May,\"%s\",%s,%s,0.200,1.440\n", testProjectNameAlt, testProjectIDAlt, testContactPersonAlt)

	chargeback := NewChargeBack(
		[]Transfer{
			NewTransfer(
				testProjectName,
				testProjectID,
				testMonth,
				testContactPerson,
				testAmount,
			),
			NewTransfer(
				testProjectNameAlt,
				testProjectIDAlt,
				testMonth,
				testContactPersonAlt,
				testAmountAlt,
			),
		},
		testMargin,
		testCurrency,
	)

	// when
	accountingDoc := GenerateAccountingDocument(chargeback)

	if expected != accountingDoc {
		t.Errorf("%v should be %v", accountingDoc, expected)
	}

	// then
}

func TestNewChargeBack(t *testing.T) {
	//given ....
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
	entry := NewTransfer(testProjectName, testProjectID, testMonth, testContactPerson, testAmount)

	entries := []Transfer{
		entry,
	}
	testMargin := 0.2
	testCurrency := "$"

	// when ....
	chargeback := NewChargeBack(entries, testMargin, testCurrency)

	// then
	for i, inputEntry := range chargeback.transfers {
		if !cmp.Equal(entries[i], inputEntry) {
			t.Fail()
		}
	}
}

func TestNewTransfer(t *testing.T) {
	testProjectName := "Doe Company"
	testProjectID := "12345"
	testMonth := time.Now()
	testContactPerson := "John Doe"
	testAmount := 1.44
	entry := NewTransfer(testProjectName, testProjectID, testMonth, testContactPerson, testAmount)

	if entry.Month != testMonth ||
		entry.ProjectID != testProjectID ||
		entry.ContactPerson != testContactPerson ||
		entry.Amount != testAmount {
		t.Fail()
	}

}

// TODO: test pnce more with DecimalCommas on
// Consider: Passing the decimal operator explicitly may be a better solution, as it avoids different control flow branches
func TestFormatRow(t *testing.T) {
	// given
	monthLayout := "2006-Jan"
	monthString := "2019-May"
	testmonth, err := time.Parse(monthLayout, monthString)
	if err != nil {
		t.Errorf("Error in test setup - Can't parse testmonth")
	}

	testpos := 0
	testName := "testName"
	testID := "123"
	testContactPerson := "John Doe"
	testMargin := 0.2
	testAmount := 1.23456789
	expected := []string{"1", monthString, testName, testID, testContactPerson, "0.200", "1.235"}

	row := formatRow(testpos, testmonth, testName, testID, testContactPerson, testMargin, testAmount, false)

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
