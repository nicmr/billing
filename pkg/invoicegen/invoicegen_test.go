package invoicegen

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewGeneratorInput(t *testing.T) {
	testMonth := time.Now()
	testProjectID := "12345"
	testContactPerson := "John Doe"
	testAmount := 1.2
	testMargin := 0.2
	testTotal := 1.44
	entry := NewGeneratorEntry(testMonth, testProjectID, testContactPerson, testAmount, testMargin, testTotal)

	entries := []GeneratorEntry{
		entry,
	}

	input := NewGeneratorInput(entries)

	for i, inputEntry := range input.Entries {
		if !cmp.Equal(entries[i], inputEntry) {
			t.Fail()
		}
	}
}

func TestNewGeneratorEntry(t *testing.T) {
	testMonth := time.Now()
	testProjectID := "12345"
	testContactPerson := "John Doe"
	testAmount := 1.2
	testMargin := 0.2
	testTotal := 1.44
	entry := NewGeneratorEntry(testMonth, testProjectID, testContactPerson, testAmount, testMargin, testTotal)

	if entry.Month != testMonth ||
		entry.ProjectID != testProjectID ||
		entry.ContactPerson != testContactPerson ||
		entry.Amount != testAmount ||
		entry.Margin != testMargin ||
		entry.Total != testTotal {
		t.Fail()
	}

}
