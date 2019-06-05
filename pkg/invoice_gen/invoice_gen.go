package invoice_gen

import "time"

func DefaultFormat(input GeneratorInput) string {
	return CSV(input)
}

// GeneratorInput contains all information required by the functions of this package
type GeneratorInput struct {
	Entries []GeneratorEntry
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
