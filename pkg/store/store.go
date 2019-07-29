// Package store provides functions to arbitrary file contents on AWS S3
package store

import (
	"fmt"
	"io/ioutil"
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
	wg      sync.WaitGroup
	Outputs []UploadResult
}

// UploadResult contains information about the success and metadata
// about an single upload in an UploadGroup
type UploadResult struct {
	S3Output chan *s3manager.UploadOutput
	Err      chan error
}

// Wait blocks execution and waits for all Uploads of the UploadGroup to finish
func (group *UploadGroup) Wait() {
	group.wg.Wait()
}

// Upload starts a new file upload goroutine for the UploadGroup with the specified parameters
// It returns an error channel it will write any encountered errors to.
// If no errors are encountered, it will write nil to the channel.
func (group *UploadGroup) Upload(contents string, bucket string, filename string, month time.Time) {
	group.wg.Add(1)
	uploadResult := UploadResult{
		make(chan *s3manager.UploadOutput, 1),
		make(chan error, 1),
	}

	go func(outc chan *s3manager.UploadOutput, ec chan error) {
		defer group.wg.Done()
		output, err := Upload(contents, bucket, filename, month)
		ec <- err
		outc <- output
	}(uploadResult.S3Output, uploadResult.Err)

	group.Outputs = append(group.Outputs, uploadResult)
}

// LocalFile creates a local file with the specified content, using the same naming scheme as the Uploads to S3
func LocalFile(contents string, dirPath string, filename string, month time.Time) error {
	filenameKeyScheme := s3KeyScheme(month, filename)

	err := ioutil.WriteFile(dirPath+filenameKeyScheme, []byte(contents), 0644)
	if err != nil {
		return err
	}
	return nil
}

// Upload uploads bytes from `reader` to S3 bucket `bucket`.
// The created objects name will be created according to store.s3KeyScheme
func Upload(contents string, bucket string, filename string, month time.Time) (*(s3manager.UploadOutput), error) {

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
