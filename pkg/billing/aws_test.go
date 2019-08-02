package billing

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func TestAWS(t *testing.T) {
	// no parameters will look for credentials file at awscli default location
	credsFile := credentials.NewSharedCredentials("", "")
	credsEnv := credentials.NewEnvCredentials()
	if credsFile == nil || credsEnv == nil {
		log.Println("AWS unexpectedly returned a nil pointer")
		t.SkipNow()
	}
	_, errFileCreds := credsFile.Get()
	_, errCredsEnv := credsEnv.Get()
	if errFileCreds != nil && errCredsEnv != nil {
		log.Println("Skipping Test because Creds could be retrieved neither from file or env.")
		t.SkipNow()
	}

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

func TestSplitIntoBounds(t *testing.T) {
	const iso8601 = "2006-01-02"
	{
		// given 31 day month
		samplemonth := "2019-03-22"
		month, err := time.Parse(iso8601, samplemonth)
		if err != nil {
			t.Errorf("Error in TestSplitIntoBounds setup: can't parse date: %v", err)
		}

		// when
		first, last := splitIntoBounds(month)

		// then
		if first != "2019-03-01" ||
			last != "2019-04-01" {
			t.Errorf("Wrong date bounds for 31 day month")
		}
	}

	{
		// given 30 day month
		samplemonth := "2019-04-23"
		month, err := time.Parse(iso8601, samplemonth)
		if err != nil {
			t.Errorf("Error in TestSplitIntoBounds setup: can't parse date: %v", err)
		}

		// when
		first, last := splitIntoBounds(month)

		// then
		if first != "2019-04-01" ||
			last != "2019-05-01" {
			t.Errorf("Wrong date bounds for 30 day month")
		}
	}

	{
		// given 28 day month
		samplemonth := "2019-02-21"
		month, err := time.Parse(iso8601, samplemonth)
		if err != nil {
			t.Errorf("Error in TestSplitIntoBounds setup: can't parse date: %v", err)
		}

		// when
		first, last := splitIntoBounds(month)

		// then
		if first != "2019-02-01" ||
			last != "2019-03-01" {
			t.Errorf("Wrong date bounds for 28 day month")
		}
	}

	{
		// given 29 day month
		samplemonth := "2016-02-21"
		month, err := time.Parse(iso8601, samplemonth)
		if err != nil {
			t.Errorf("Error in TestSplitIntoBounds setup: can't parse date: %v", err)
		}

		// when
		first, last := splitIntoBounds(month)

		// then
		if first != "2016-02-01" ||
			last != "2016-03-01" {
			t.Errorf("Wrong date bounds for 29 day month")
		}
	}

}
