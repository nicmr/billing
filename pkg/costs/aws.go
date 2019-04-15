package costs

import (
	"log"
	"time"

	"github.com/Altemista/altemista-billing/pkg/csv"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var (
	// Session is safe for concurrent use after initialization,
	// as it will not be mutated by the SDK after creation
	awsSess = createSessionOrFatal()
)

func createSessionOrFatal() *(session.Session) {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Unable to initialize aws session: ", err)
	}
	return sess
}

// costsBetween returns the a GetCostAndUsageOutput containing the costs created between `start` and `end`.
// Start and end should be strings of the form "YYYY-MM-DD".
// This date range is left-inclusive and right-exclusive.
func costexplorerCall(costexpl *(costexplorer.CostExplorer), start string, end string, metrics []*string) (*costexplorer.GetCostAndUsageOutput, error) {
	// truestring := "true"

	// prepare a GetCostAndUsageInput struct for the request
	input := (&costexplorer.GetCostAndUsageInput{}).
		SetTimePeriod((&costexplorer.DateInterval{}).
			SetStart(start).
			SetEnd(end)).
		SetGranularity("MONTHLY").
		// SetFilter((&costexplorer.Expression{}).
		// 	SetTags((&costexplorer.TagValues{}).
		// 		SetKey("isUserResource").
		// 		SetValues([]*string{&truestring}))).
		SetGroupBy([]*costexplorer.GroupDefinition{(&costexplorer.GroupDefinition{}).
			// Key can be AZ, INSTANCE_TYPE,  LEGAL_ENTITY_NAME, LINKED_ACCOUNT, OPERATION, PLATFORM,
			// PURCHASE_TYPE, SERVICE, TENANCY, and USAGE_TYPE, if type is DIMENSION.
			// It can be the name of any cost explorer tag, if type is TAG
			// SetKey("customerID").
			// Type needs to be either DIMENSION or TAG
			// SetType("TAG")}).
			SetKey("INSTANCE_TYPE").
			SetType("DIMENSION")}).
		SetMetrics(metrics)

	output, err := costexpl.GetCostAndUsage(input)
	if err != nil {
		return nil, err
	}
	log.Println(output.String())
	return output, nil
}

func maxGroupLen(arr []*costexplorer.ResultByTime) int {
	max := 0
	for _, e := range arr {
		if max < len(e.Groups) {
			max = len(e.Groups)
		}
	}
	return max
}

// costsBetweenAWS calls costexplorer after adding package-level variables as parameters,
// then timestamps the result, generates cooresponding csv and returns it as an APICallResult
func costsBetweenAWS(month string) (APICallResult, error) {
	start, end, err := splitIfValid(month)
	if err != nil {
		return APICallResult{}, err
	}

	amortizedCost := "AmortizedCost"
	metrics := []*string{&amortizedCost}

	output, err := costexplorerCall(costexplorer.New(awsSess), start, end, metrics)
	if err != nil {
		return APICallResult{}, err
	}

	csvEntries := make([]csv.Entry, maxGroupLen(output.ResultsByTime))

	iso8601Short := "2006-01"
	desiredFormat := "2006-Jan"

	// Retrieve the required information for csvEntries from the output.
	// this implementation only works for a single month
	element := output.ResultsByTime[0]
	m, err := time.Parse(iso8601Short, month)
	if err != nil {
		return APICallResult{}, err
	}
	monthStr := m.Format(desiredFormat)
	for i, group := range element.Groups {
		csvEntries[i] = csv.Entry{
			Month:         monthStr,
			ProjectID:     "Not yet implemented",
			ContactPerson: "Not yet implemented",
			Amount:        *group.Metrics[amortizedCost].Amount,
		}
	}

	result := APICallResult{
		Timestamp:      time.Now(),
		Response:       output.String(),
		CsvFileContent: csv.Marshal(csvEntries),
	}

	return result, nil
}

// splitIfValid checks if month is a iso 8601 conforming string,
// then splits it into the first day of the month and the first day fo the following month
// These two tasks are combined in one function because it is more efficient, validating it is a side effect of splitting it
func splitIfValid(month string) (string, string, error) {
	const iso8601 = "2006-01-02"
	startstr := month + "-01"

	targetMonth, err := time.Parse(iso8601, startstr)
	if err != nil {
		return "", "", err
	}
	nextmonth := targetMonth.AddDate(0, 1, 0)
	y, m, _ := nextmonth.Date()
	end := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)

	endstr := end.Format(iso8601)
	log.Println(startstr, endstr)
	return startstr, endstr, nil
}
