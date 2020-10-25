package go_sqs_extended

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SendMessage creates and inserts a record into the given SQS queue
// If large message processing is enabled, the contents of the message are inserted to S3 as a new object.
// Then, a reference to the object will be sent as the SQS record instead.
// If large message processing is disabled, the contents of the message are sent to SQS as normal.
func (esc *ExtendedSQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.SendMessage(input)
	}
	if !esc.cfg.AlwaysSendThroughS3 && !esc.messageIsLarge(input) {
		return esc.SQS.SendMessage(input)
	}
	id, err := esc.s3c.Client.WriteBigMessage(*input.MessageBody)
	if err != nil {
		return &sqs.SendMessageOutput{},
			fmt.Errorf("message could not be uploaded to S3, reason: %v. nothing was sent to SQS", err)
	}
	asBytes, err := json.Marshal(&ExtendedQueueMessage{
		S3BucketName: esc.s3c.BucketName,
		S3Key:        id,
	})
	if err != nil {
		return &sqs.SendMessageOutput{}, fmt.Errorf("couldn't parse SendMessageInput message, %v", err)
	}
	// shallow copy to avoid side effects (modifying input object)
	duplicate := *input
	duplicate.MessageBody = aws.String(string(asBytes))
	return esc.SQS.SendMessage(&duplicate)
}

func (esc *ExtendedSQS) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.SendMessageBatch(input)
	}
	if esc.cfg.AlwaysSendThroughS3 || esc.batchMessageIsLarge(input) {
		return esc.SQS.SendMessageBatch(input)
	}
	// id, err := esc.s3c.Client.WriteBigMessage(*input.MessageBody)
	return esc.SQS.SendMessageBatch(input)
}

func (esc *ExtendedSQS) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return nil, nil
}

func (esc *ExtendedSQS) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

func (esc *ExtendedSQS) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	return nil, nil
}

func (esc *ExtendedSQS) ChangeMessageVisibility(input *sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	return nil, nil
}

func (esc *ExtendedSQS) ChangeMessageVisibilityBatch(input *sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	return nil, nil
}
