// Package store provides functions to arbitrary file contents on AWS S3
package store

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	// Session is safe for concurrent use after initialization,
	// as it will not be mutated by the SDK after creation
	awsSess *session.Session
	// Uploader is also safe for concurrent use
	uploader *s3manager.Uploader
)

func init() {
	// initialize aws session
	var err error
	awsSess, err = session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		log.Fatal("Unable to initialize aws session for package store", err)
	}

	// initialize uploader
	uploader = s3manager.NewUploader(awsSess)
}

// UploadGroup facilitates uploading multiple files concurrently, then waiting for all of them to finish.
// It is a safe abstraction around sync.WaitGroup
type UploadGroup struct {
	wg sync.WaitGroup
}

// Wait blocks execution and waits for all Uploads of the UploadGroup to finish
func (group *UploadGroup) Wait() {
	group.wg.Wait()
}

// Upload starts a new file upload goroutine for the UploadGroup with the specified parameters
// It returns an error channel it will write any encountered errors to.
// If no errors are encountered, it will write nil to the channel.
func (group *UploadGroup) Upload(contents string, bucket string, filename string, fileExtension string, month time.Time) chan error {
	group.wg.Add(1)
	errchan := make(chan error, 1)

	go func(ec chan error) {
		defer group.wg.Done()
		_, err := Upload(contents, bucket, filename, fileExtension, month)
		ec <- err
	}(errchan)

	return errchan
}

// Upload uploads bytes from `reader` to S3 bucket `bucket`.
// The created objects name will be created according to store.s3KeyScheme
func Upload(contents string, bucket string, filename string, fileExtension string, month time.Time) (*(s3manager.UploadOutput), error) {

	reader := strings.NewReader(contents)

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
	parentprefix := "invoices"
	yearprefix := month.Format("2006")
	monthprefix := month.Format("Jan")

	monthformat := "2006-01"
	monthstr := month.Format(monthformat)

	key := fmt.Sprintf("%s/%s/%s/%s_%s_%s.%s",
		parentprefix, yearprefix, monthprefix,
		filename, monthstr, time.Now().Format("2006-01-02-15:04:05"),
		"csv")

	return key
}
