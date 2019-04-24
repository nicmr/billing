package s3store

import (
	"fmt"
	"io"
	"log"
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
// The created objects name will be {fileName}{timestamp}?{fileExtension}
func Upload(reader io.Reader, bucket string, fileName string, fileExtension string, includeTimestamp bool) (*(s3manager.UploadOutput), error) {
	uploader := s3manager.NewUploader(awsSess)
	fullKey := fileName
	if includeTimestamp {
		now := time.Now().Format("_2006-01-02_15:04:05")
		fullKey += now
	}
	fullKey += fileExtension

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fullKey),
		Body:   reader,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file, %v", err)
	}

	log.Println("Uploaded as fileName: ", fullKey)

	return result, nil
}
