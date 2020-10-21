package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type extendedSQSClient struct {
	sqsiface.SQSAPI
	Configuration *ExtendedConfiguration
}

type ExtendedConfiguration struct {

}

func NewExtended(sqs *sqs.SQS, options *ExtendedConfiguration) sqsiface.SQSAPI {
	return &extendedSQSClient{
		SQSAPI:             sqs,
		Configuration:   options,
	}
}

func (esc *extendedSQSClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) ChangeMessageVisibility(input *sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	return nil, nil
}

func (esc *extendedSQSClient) ChangeMessageVisibilityBatch(input *sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	return nil, nil
}
