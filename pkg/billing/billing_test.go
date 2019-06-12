package billing

import (
	"testing"
	"time"
)

func TestCalculateCosts(t *testing.T) {

	// create a mocked provider
	testingProvider := CloudProvider{
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
		t.Errorf("error parsing date string %v", err)
	}

	// Call the function we want to test with the created parameters
	_, err = CalculateCosts(testingProvider, month, 0.0)
	if err != nil {
		t.Errorf("Failed calculating costs: %v", err)
	}

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
		t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, totals[0])
	}
	if totals[1] != 0.05*1.2 {
		t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, totals[1])
	}

	// Tests with margin = 0
	testMargin = 0.0
	totals = applyMargin(testEntries, testMargin)

	if totals[0] != 1000.0 {
		t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, totals[0])
	}
	if totals[1] != 0.05 {
		t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, totals[1])
	}

}
