package go_sqs_extended

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hashicorp/go-multierror"
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
	asMessage, err := esc.s3c.Client.CreateBigMessage(input.MessageBody)
	if err != nil {
		return &sqs.SendMessageOutput{},
			fmt.Errorf("message could not be uploaded to S3, nothing was sent to SQS: %w", err)
	}
	asBytes, err := json.Marshal(&asMessage)
	if err != nil {
		return &sqs.SendMessageOutput{}, fmt.Errorf("couldn't parse SendMessageInput message: %w", err)
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
	extendedInput := *input
	extendedInput.AttributeNames = esc.enforceSingleReserved(input.AttributeNames)
	result, err := esc.SQS.ReceiveMessage(&extendedInput)

	if !esc.s3c.Configured || err != nil {
		return result, err
	}

	var merr error
	extracted := make([]*sqs.Message, len(result.Messages))
	for i, message := range result.Messages {
		if m, err := esc.createExtractedMessage(message); err != nil {
			merr = multierror.Append(merr, err)
		} else {
			extracted[i] = m
		}
	}

	return &sqs.ReceiveMessageOutput{Messages: extracted}, merr
}

func (esc *ExtendedSQS) createExtractedMessage(message *sqs.Message) (*sqs.Message, error) {
	var extendedBody = new(extendedQueueMessage)
	err := json.Unmarshal([]byte(*message.Body), extendedBody)
	if err != nil {
		return message, fmt.Errorf("couldn't extract message body: %w", err)
	}
	contents, err := esc.s3c.Client.ExtractBigMessage(extendedBody)
	if err != nil {
		return message, fmt.Errorf("failed to extract message %s from s3, result not modified: %w",
			*message.MessageId, err)
	}
	hashed := md5.Sum([]byte(contents))
	handle := esc.embedS3PointerInReceiptHandle(
		aws.StringValue(message.ReceiptHandle), extendedBody,
	)
	return &sqs.Message{
		Attributes:             message.Attributes,
		Body:                   aws.String(contents),
		MD5OfBody:              aws.String(fmt.Sprintf("%x", hashed)),
		MD5OfMessageAttributes: message.MD5OfMessageAttributes,
		MessageAttributes:      message.MessageAttributes,
		MessageId:              message.MessageId,
		ReceiptHandle:          aws.String(handle),
	}, err
}

func (esc *ExtendedSQS) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.DeleteMessage(input)
	}
	return nil, nil
}

func (esc *ExtendedSQS) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.DeleteMessageBatch(input)
	}
	return nil, nil
}

func (esc *ExtendedSQS) ChangeMessageVisibility(input *sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.ChangeMessageVisibility(input)
	}
	return nil, nil
}

func (esc *ExtendedSQS) ChangeMessageVisibilityBatch(input *sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	if !esc.s3c.Configured {
		return esc.SQS.ChangeMessageVisibilityBatch(input)
	}
	return nil, nil
}
