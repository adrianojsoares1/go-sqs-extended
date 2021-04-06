package go_sqs_extended

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type Mock_SendMessage struct {
	sqsiface.SQSAPI
	OutputFn func(*sqs.SendMessageInput) *sqs.SendMessageOutput
}

type Mock_RecieveMessage struct {
	sqsiface.SQSAPI
	OutputFn func(*sqs.ReceiveMessageInput) *sqs.ReceiveMessageOutput
}

type Mock_SendMessageBatch struct {
	sqsiface.SQSAPI
	OutputFn func(*sqs.SendMessageBatchInput) *sqs.SendMessageBatchOutput
}

type Mock_DeleteMessage struct {
	sqsiface.SQSAPI
	OutputFn func(*sqs.DeleteMessageInput) *sqs.DeleteMessageOutput
}

type Mock_BigS3MessageClient struct {
	BigMessageS3Client
	CreateOutput *extendedQueueMessage
	ExtractOutput string
}

func (b3mc *Mock_BigS3MessageClient) CreateBigMessage(message *string) (*extendedQueueMessage, error) {
	return b3mc.CreateOutput, nil
}

func (b3mc *Mock_BigS3MessageClient) ExtractBigMessage(message *extendedQueueMessage) (string, error) {
	return b3mc.ExtractOutput, nil
}

func (msm *Mock_SendMessage) SendMessage(smi *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return msm.OutputFn(smi), nil
}

func (mrm *Mock_RecieveMessage) ReceiveMessage(rmi *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return mrm.OutputFn(rmi), nil
}

func TestExtendedSQS_SendMessage(t *testing.T) {}
