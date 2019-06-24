package store

import (
	"log"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	_, err := Upload("test", "altemista-billing-travis", "test/invoice", ".csv", time.Now())
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
		t.FailNow()
	}
}
