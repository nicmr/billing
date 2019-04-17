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
			// // Key can be AZ, INSTANCE_TYPE,  LEGAL_ENTITY_NAME, LINKED_ACCOUNT, OPERATION, PLATFORM,
			// // PURCHASE_TYPE, SERVICE, TENANCY, and USAGE_TYPE, if type is DIMENSION.
			// // It can be the name of any cost explorer tag, if type is TAG
			// SetKey("customerID").
			// // Type needs to be either DIMENSION or TAG
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

// costsBetweenAWS calls costexplorer after adding package-level variables as parameters,
// then timestamps the result, generates cooresponding csv and returns it as an APICallResult
func costsMonthlyAWS(month time.Time) (APICallResult, error) {

	amortizedCost := "AmortizedCost"
	metrics := []*string{&amortizedCost}

	start, end := splitIntoBounds(month)

	output, err := costexplorerCall(costexplorer.New(awsSess), start, end, metrics)
	if err != nil {
		return APICallResult{}, err
	}

	// reserve space for the queried month
	csvEntries := make([]csv.Entry, len(output.ResultsByTime[0].Groups))

	desiredFormat := "2006-Jan"

	monthStr := month.Format(desiredFormat)

	// Retrieve the required information for csvEntries from the output.
	// As we queried only for a single month, we don't have to iterate and simply look at [0]
	element := output.ResultsByTime[0]
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

// splitIntoBounds splits month into the first day of the month and the first day fo the following month
// These two tasks are combined in one function because it is more efficient, validating it is a side effect of splitting it
func splitIntoBounds(month time.Time) (string, string) {
	const iso8601 = "2006-01-02"

	y, m, _ := month.Date()
	firstOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	start := firstOfMonth.Format(iso8601)
	end := firstOfNextMonth.Format(iso8601)

	return start, end
}
