package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {

	// TODO: Make sure application reads secrets from Credentials File (dev) or IAM Role(production)

	// Create new session and costexplorer client
	sess, err := session.NewSession()
	if err != nil {
		log.Println("Unable to create session", err)
	}
	log.Println("region: ", *(sess.Config.Region))

	svc := costexplorer.New(sess)

	// TODO: Move date parameters to flags

	output, err := costsBetween(svc, "2019-03-31", "2019-04-02")

	//yyyy-MM-ddThh:mm:ssZ
	// output, err := costsBetween(svc, "2019-03-31T00:00:00Z", "2019-04-02T00:00:00Z")

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
	}

	log.Println(output.String())

	// TODO: Parse `output` and generate desired output
}

// costsBetween returns the a GetCostAndUsageOutput containing the costs created between `start` and `end`.
// Start and end should be strings of the form "YYYY-MM-DD".
// This date range is left-inclusive and right-exclusive.
func costsBetween(costexpl *(costexplorer.CostExplorer), start string, end string) (*costexplorer.GetCostAndUsageOutput, error) {
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
