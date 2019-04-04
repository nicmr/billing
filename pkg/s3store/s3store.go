package s3store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
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

// Upload uploads `reader` to the program's associated s3 bucket and saves it under
// key, where key is a path-style object key for the buckett.
// TODO: Specify s3 bucket as additional parameter or as config loaded from a config file?
func Upload(reader io.Reader, key string) (*(s3manager.UploadOutput), error) {
	uploader := s3manager.NewUploader(awsSess)
	myBucket := "altemista-billing"

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file, %v", err)
	}
	return result, nil
}
