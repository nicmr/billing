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
	output, err := costsBetween(svc, "2019-03-20", "2019-03-29")
	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
	}

	log.Println(output.String())

	// TODO: Parse `output` and generate desired output
}

// costsBetween is an abstraction around the creation of a GetCostAndUsageInput
func costsBetween(costexpl *(costexplorer.CostExplorer), start string, end string) (*costexplorer.GetCostAndUsageOutput, error) {
	// input := costexplorer.GetCostAndUsageInput{}
	// interval := (&costexplorer.DateInterval{}).SetStart(start).SetEnd(end)
	// input.SetTimePeriod(interval)
	truestring := "true"

	input := (&costexplorer.GetCostAndUsageInput{}).
		SetTimePeriod((&costexplorer.DateInterval{}).
			SetStart(start).
			SetEnd(end)).
		SetFilter((&costexplorer.Expression{}).
			SetTags((&costexplorer.TagValues{}).
				SetKey("isUserResource").
				SetValues([]*string{&truestring}))).
		SetGroupBy([]*costexplorer.GroupDefinition{(&costexplorer.GroupDefinition{}).
			SetKey("customerID")})

	output, err := costexpl.GetCostAndUsage(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}
