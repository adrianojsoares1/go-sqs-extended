package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (esc *ExtendedSQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if !esc.Configuration.S3Configuration.isConfigured() {
		return esc.SQS.SendMessage(input)
	}
	if esc.Configuration.AlwaysSendThroughS3 || esc.messageIsLarge(input) {
	}
	return esc.SQS.SendMessage(input)
}

func (esc *ExtendedSQS) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if !esc.Configuration.S3Configuration.isConfigured() {
		return esc.SQS.SendMessageBatch(input)
	}
	if esc.Configuration.AlwaysSendThroughS3 || esc.batchMessageIsLarge(input) {

	}
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
