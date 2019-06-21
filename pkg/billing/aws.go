package billing

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

var (
	// Session is safe for concurrent use after initialization,
	// as it will not be mutated by the SDK after creation
	awsSess = createSessionOrFatal()
)

// costsBetweenAWS calls costexplorer after adding package-level variables as parameters,
// then timestamps the result, generates cooresponding csv and returns it as an APICallResult
func costsMonthlyAWS(month time.Time) (apiCallResult, error) {

	amortizedCost := "AmortizedCost"
	metrics := []*string{&amortizedCost}

	start, end := splitIntoBounds(month)

	output, err := costexplorerCall(costexplorer.New(awsSess), start, end, metrics)
	if err != nil {
		return apiCallResult{}, err
	}

	// reserve space for the queried month
	entries := make([]apiCallResultEntry, len(output.ResultsByTime[0].Groups))

	// Retrieve the required information for csvEntries from the output.
	// As we queried only for a single month, we don't have to iterate and simply look at [0]
	element := output.ResultsByTime[0]
	for i, group := range element.Groups {
		amount, err := strconv.ParseFloat(*group.Metrics[amortizedCost].Amount, 64)
		if err != nil {
			log.Println("Unable to decode AWS cost amount as float64")
			return apiCallResult{}, err
		}

		entries[i] = apiCallResultEntry{
			ProjectID: strings.Replace(*group.Keys[0], "project-number$", "", 1),
			Amount:    amount,
			Currency:  *(group.Metrics[amortizedCost].Unit),
		}
	}
	result := apiCallResult{
		Timestamp:      time.Now(),
		ResponseString: output.String(),
		Entries:        entries,
	}

	return result, nil
}

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
// It is safer to use the wrapping function costsMonthlyAWS instead to get information about a single month
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
			SetKey("project-number").
			// SetKey("INSTANCE_TYPE").
			// // Type needs to be either DIMENSION or TAG
			SetType("TAG")}).
		// SetType("DIMENSION")}).
		SetMetrics(metrics)

	output, err := costexpl.GetCostAndUsage(input)
	if err != nil {
		return nil, err
	}
	log.Println(output.String())
	return output, nil
}

// splitIntoBounds splits month into the first day of the month and the first day fo the following month
func splitIntoBounds(month time.Time) (string, string) {
	const iso8601 = "2006-01-02"

	y, m, _ := month.Date()
	firstOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	start := firstOfMonth.Format(iso8601)
	end := firstOfNextMonth.Format(iso8601)

	return start, end
}
