package costs

import (
	"log"
	"testing"
)

func TestAWS(t *testing.T) {
	apicall := AWS()
	_, err := apicall("2019-04")
	if err != nil {
		log.Println("costBetween call failed: ", err)
		t.FailNow()
	}
}
