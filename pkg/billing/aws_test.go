package billing

import (
	"log"
	"testing"
	"time"
)

func TestAWS(t *testing.T) {
	const iso8601 = "2006-01-02"
	provider := AWS()
	month := "2019-04-01"
	parsedMonth, err := time.Parse(iso8601, month)
	if err != nil {
		log.Println("error parsing date string", err) // should never happen
		t.FailNow()
	}
	_, err = provider.apicall(parsedMonth)
	if err != nil {
		log.Println("costBetween call failed: ", err)
		t.FailNow()
	}
}
