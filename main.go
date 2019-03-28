package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {

	// TODO: Create a custom session config
	// Ref: https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-central-1")})
	if err != nil {
		log.Println("Unable to create session", err)
	}

	svc := costexplorer.New(sess)

	// TODO: Modify the input struct to apply filters
	// Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/costexplorer/#GetCostAndUsageInput
	input := costexplorer.GetCostAndUsageInput{}
	output, err := svc.GetCostAndUsage(&input)
	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
	}

	log.Println(output.String())

	// TODO: Generate desired output data from `output`
}
