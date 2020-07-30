package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/pfeilbr/aws-sam-golang-example/job"
	"github.com/pfeilbr/aws-sam-golang-example/lambdautils"
)

// Worker consumes the messages and executes the job.
func Worker(ctx context.Context, event events.SQSEvent, svc sqsiface.SQSAPI) error {
	var err error

	for _, message := range event.Records {

		// Do the "work" here.
		if ctx, err = job.Do(ctx, message); err != nil {
			return err
		}

		// Delete the message from SQS so it is not processed again.
		lambdautils.DeleteMessage(svc, message.ReceiptHandle)
	}

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	return Worker(ctx, event, svc)
}

func main() {
	lambda.Start(handler)
}
