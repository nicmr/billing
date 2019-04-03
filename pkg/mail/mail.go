package mail

import (
// "github.com/aws/aws-sdk-go/aws/session"
// "github.com/aws/aws-sdk-go/service/ses"
)

// Why SES and not SNS?
// SNS requires the definition of custom endpoints on the client side
// As we don't have a client, communication has to take place exclusively by mail

// Concept:
// People can be added as subscribers to an S3 bucket
// When the analysis of billing data has launched,
// the service imports the current subscribers and sends a mail to each
//

// Subscribe adds the specified E-Mail adress as a new subscriber to the connected s3 bucket
func Subscribe(adress string) error {
	return nil
}

// ImportSubscribers imports the current subscribers from a S3 bucket
func ImportSubscribers() error {
	return nil
}

// SendEMails sends emails to all subcribers
func SendEMails() error {
	return nil
}
