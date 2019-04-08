package costs

import (
	"log"
	"time"

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
func costexplorerCall(costexpl *(costexplorer.CostExplorer), start string, end string) (*costexplorer.GetCostAndUsageOutput, error) {
	// truestring := "true"
	metrics := "AmortizedCost"

	// prepare a GetCostAndUsageInput struct for the request
	input := (&costexplorer.GetCostAndUsageInput{}).
		SetTimePeriod((&costexplorer.DateInterval{}).
			SetStart(start).
			SetEnd(end)).
		SetGranularity("DAILY").
		// SetFilter((&costexplorer.Expression{}).
		// 	SetTags((&costexplorer.TagValues{}).
		// 		SetKey("isUserResource").
		// 		SetValues([]*string{&truestring}))).
		SetGroupBy([]*costexplorer.GroupDefinition{(&costexplorer.GroupDefinition{}).
			SetKey("customerID").
			SetType("TAG")}).
		SetMetrics([]*string{&metrics})

	output, err := costexpl.GetCostAndUsage(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// CostsBetweenAWS returns
// It adds package-level variables as parameters, forwards the function call and adds a timestamp
func costsBetweenAWS(start string, end string) (APICallResult, error) {

	output, err := costexplorerCall(costexplorer.New(awsSess), start, end)
	if err != nil {
		return APICallResult{}, err
	}

	result := APICallResult{
		Timestamp: time.Now(),
		Response:  output.String(),
	}

	return result, nil
}
