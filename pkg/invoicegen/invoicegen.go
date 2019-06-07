package invoicegen

import "time"

// DefaultFormat selects the recommended default invoice format generator.
// Useful if you require no specific format and would like to not have to make changes inn your code once the recommended deafult changes.
// If you require a specific format, such as CSV, please use the corresponding function
func DefaultFormat(convertable ConvertableToGeneratorInput) string {
	return CSV(convertable)
}

// GeneratorInput contains all information required by the functions of this package
type GeneratorInput struct {
	Entries []GeneratorEntry
}

// NewGeneratorInput returns a GeneratorInput to be passed to the Different methods of the package
// This is the preferred method of instantiating a struct of this type
func NewGeneratorInput(entries []GeneratorEntry) GeneratorInput {
	return GeneratorInput{entries}
}

// NewGeneratorEntry returns a GeneratorEntry to be used with a GeneratorInput struct.
// This is the preferred method of instantiating a struct of this type
func NewGeneratorEntry(month time.Time, projectID string, contactPerson string, amount float64, margin float64, total float64) GeneratorEntry {
	return GeneratorEntry{
		Month:         month,
		ProjectID:     projectID,
		ContactPerson: contactPerson,
		Amount:        amount,
		Margin:        margin,
		Total:         total,
	}
}

// GeneratorEntry represents a single entry of specific information in the GeneratorInput struct
// Usually represents a row in 2D based generation targets
type GeneratorEntry struct {
	Month         time.Time
	ProjectID     string
	ContactPerson string
	Amount        float64
	Margin        float64
	Total         float64
}

// ConvertableToGeneratorInput allows you to define a function for converting your structs to `GeneratorInput`s
// This allows them to be passed directly to the functions in the public interface of this package.
type ConvertableToGeneratorInput interface {
	ToInvoiceGenInput() GeneratorInput
}
