package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"testing"
)

func TestExtendedSQS_messageIsLarge(t *testing.T) {
	cases := []struct {
		Name string
		ThresholdKb int64
		Input *sqs.SendMessageInput
		ExpectedLarge bool
	}{
		{
			Name: "Message Is Large",
			ThresholdKb: DefaultLargeMessageSize,
			Input: &sqs.SendMessageInput{
				MessageAttributes:       map[string]*sqs.MessageAttributeValue {
					"four" : { BinaryValue: []byte("this is a message of size 28") },
				},
				MessageBody:             aws.String("this is a message of size 28"),
			},
			ExpectedLarge: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			//v := big_message_aws.NewS3Client()
		})
	}
}

func TestExtendedSQS_batchMessageIsLarge(t *testing.T) {
	cases := []struct {
		Name string
		Input *sqs.SendMessageInput
		ExpectedSize int64
	}{
		{

		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {

		})
	}
}
