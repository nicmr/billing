package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {

	// TODO: Make sure the program reads the config stored at $HOME/.aws/config
	// (config should be separate from code)
	// Ref: https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
	sess, err := session.NewSession()
	if err != nil {
		log.Println("Unable to create session", err)
	}
	log.Println("region: ", *(sess.Config.Region))

	svc := costexplorer.New(sess)

	// TODO: Modify the input struct to apply filters
	// Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/costexplorer/#GetCostAndUsageInput

	output, err := costsBetween(svc, "2019-03-20", "2019-03-29")
	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
	}

	log.Println(output.String())

	// TODO: Parse `output` and generate desired logic
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
