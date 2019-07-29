package store

import (
	"log"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestUpload(t *testing.T) {
	_, err := Upload("test", "altemista-billing-travis", "test/invoice", time.Now())
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
		t.FailNow()
	}
}

func TestS3KeyScheme(t *testing.T) {
	// given
	monthstr := "2019-05-01"
	iso8601 := "2006-01-02"
	month, err := time.Parse(iso8601, monthstr)
	if err != nil {
		t.Errorf("Error in test setup - Can't parse testmonth")
	}

	// when
	key := s3KeyScheme(month, "foo")

	// then
	{
		expected := "invoices/2019/May/foo_2019-05_"

		// cut off timestamp and file extension, which are always the last 23 runes of the key
		testableKey := string([]rune(key)[0 : len(key)-23])

		if !cmp.Equal(testableKey, expected) {
			t.Errorf("%v should be %v", testableKey, expected)
		}
	}

	// file extension
	{
		expected := ".csv"

		// extract file extension, which is always the last four runes of the key
		extension := string([]rune(key)[len(key)-4:])
		if !cmp.Equal(extension, expected) {
			t.Errorf("%v should be %v", extension, expected)
		}
	}

}
