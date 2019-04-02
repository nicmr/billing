package main

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func TestCostBetween(t *testing.T) {
	sess, err := session.NewSession()
	if err != nil {
		log.Println("Unable to create session: ", err)
		t.FailNow()
	}
	costexpl := costexplorer.New(sess)
	_, err = costsBetween(costexpl, "2019-03-31", "2019-04-02")
	if err != nil {
		log.Println("costBetween call failed: ", err)
		t.FailNow()
	}
}
