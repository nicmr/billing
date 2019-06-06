package costs

import (
	"log"
	"testing"
	"time"
)

func TestCostCalc(t *testing.T) {

	// create a mocked provider
	testingProvider := Provider{
		apicall: func(time.Time) (apiCallResult, error) {
			result := apiCallResult{
				Timestamp:      time.Now(),
				ResponseString: "testing",
				Entries: []apiCallResultEntry{
					apiCallResultEntry{
						ProjectID: "testing",
						Amount:    12345.0,
						Currency:  "$",
					},
				},
			}
			return result, nil
		},
	}

	// select month for testing
	const iso8601 = "2006-01-02"
	monthstring := "2019-04-01"
	month, err := time.Parse(iso8601, monthstring)
	if err != nil {
		log.Println("error parsing date string", err) // should never happen
		t.FailNow()
	}

	// Call the function we want to test with the created parameters
	CostCalc(testingProvider, month, 0.0)

}

func TestApplyMargin(t *testing.T) {
	testEntries := []apiCallResultEntry{
		apiCallResultEntry{
			ProjectID: "testing",
			Amount:    1000.0,
			Currency:  "$",
		},
		apiCallResultEntry{
			ProjectID: "testing2",
			Amount:    0.05,
			Currency:  "$",
		},
	}

	// Tests with 0 < margin < 1
	testMargin := 0.2
	totals := applyMargin(testEntries, testMargin)

	if totals[0] != 1000.0*1.2 {
		log.Printf("Failed to correctly apply margin %g to value %g", testMargin, totals[0])
	}
	if totals[1] != 0.05*1.2 {
		log.Printf("Failed to correctly apply margin %g to value %g", testMargin, totals[1])
	}

	// Tests with margin = 0
	testMargin = 0.0
	totals = applyMargin(testEntries, testMargin)

	if totals[0] != 1000.0 {
		log.Printf("Failed to correctly apply margin %g to value %g", testMargin, totals[0])
	}
	if totals[1] != 0.05 {
		log.Printf("Failed to correctly apply margin %g to value %g", testMargin, totals[1])
	}

}
