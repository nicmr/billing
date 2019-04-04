package costs

import (
	"log"
	"testing"
)

func TestCostBetween(t *testing.T) {
	awsClient := AWS{}
	_, err := awsClient.CostsBetween("2019-03-31", "2019-04-02")
	if err != nil {
		log.Println("costBetween call failed: ", err)
		t.FailNow()
	}
}
