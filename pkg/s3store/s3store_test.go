package s3store

import (
	"log"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	_, err := Upload(strings.NewReader("test"), "test/test.log")
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
		t.FailNow()
	}
}
