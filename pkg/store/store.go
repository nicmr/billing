package store

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	// Session is safe for concurrent use after initialization,
	// as it will not be mutated by the SDK after creation
	awsSess = createSessionOrFatal()
)

func createSessionOrFatal() *(session.Session) {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Unable to initialize aws session: ", err)
	}
	return sess
}

// Upload uploads bytes from `reader` to S3 bucket `bucket`.
// The created objects name will be created according to store.s3KeyScheme
func Upload(contents string, bucket string, filename string, fileExtension string, month time.Time) (*(s3manager.UploadOutput), error) {
	reader := strings.NewReader(contents)
	uploader := s3manager.NewUploader(awsSess)

	key := s3KeyScheme(month, filename)

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file, %v", err)
	}

	log.Println("Uploaded as fileName: ", key)

	return result, nil
}

func s3KeyScheme(month time.Time, filename string) string {
	iso8601 := "2006-01-02"

	yearprefix := month.Format("2006")
	monthprefix := month.Format("Jan")
	monthstr := month.Format(iso8601)

	key := fmt.Sprintf("%s/%s/%s_%s_%s.%s",
		yearprefix, monthprefix,
		filename, monthstr, time.Now().Format("2006-01-02-15:04:05"),
		"csv")

	return key
}
