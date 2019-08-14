package billing

import (
	"testing"
	"time"
)

func TestCalculateChargeBack(t *testing.T) {

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
	_, err = CalculateChargeBack(testingProvider, month, 0.0)
	if err != nil {
		t.Errorf("Failed calculating costs: %v", err)
	}

}

func TestApplyMargin(t *testing.T) {
	testEntry := apiCallResultEntry{
		ProjectID: "testing",
		Amount:    1000.0,
		Currency:  "$",
	}

	// Tests with 0 < margin
	{
		testMargin := 0.2
		total := applyMargin(testEntry, testMargin)

		if total != 1000.0*1.2 {
			t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, total)
		}
	}

	// Tests with margin = 0
	{
		testMargin := 0.0
		total := applyMargin(testEntry, testMargin)

		if total != 1000.0 {
			t.Errorf("Failed to correctly apply margin %g to value %g", testMargin, total)
		}
	}

	// TODO: Tests with margin < 0
	// Should error?
	{

	}
}
