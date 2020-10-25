package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"testing"
)

func getStringOfSize(bytes int) []byte {
	x := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		x[i] = 'x'
	}
	return x
}

func TestExtendedSQS_messageIsLarge(t *testing.T) {
	cases := []struct {
		Name string
		ThresholdKb int64
		Input *sqs.SendMessageInput
		ExpectedLarge bool
	}{
		{
			Name: "Message Is Not Large: 60 bytes",
			ThresholdKb: DefaultLargeMessageSize,
			Input: &sqs.SendMessageInput{
				MessageAttributes:       map[string]*sqs.MessageAttributeValue {
					"four" : { BinaryValue: getStringOfSize(28) },
				},
				MessageBody:             aws.String(string(getStringOfSize(28))),
			},
			ExpectedLarge: false,
		},
		{
			Name: "Message Is Large: 1052 Bytes vs 1024 Byte Threshold",
			ThresholdKb: int64(1024),
			Input: &sqs.SendMessageInput{
				MessageAttributes:       map[string]*sqs.MessageAttributeValue {
					"four" : { BinaryValue: getStringOfSize(28) },
				},
				MessageBody:             aws.String(string(getStringOfSize(1024))),
			},
			ExpectedLarge: true,
		},
		{
			Name: "Message Is Not Large: 1052 Bytes vs 2048 Byte Threshold",
			ThresholdKb: int64(2048),
			Input: &sqs.SendMessageInput{
				MessageAttributes:       map[string]*sqs.MessageAttributeValue {
					"four" : { BinaryValue: getStringOfSize(28) },
				},
				MessageBody:             aws.String(string(getStringOfSize(1024))),
			},
			ExpectedLarge: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			v := (&ExtendedSQS{
				cfg: &extendedConfigurationGhost{
					LargeMessageThreshold: tc.ThresholdKb,
				},
			}).messageIsLarge(tc.Input)
			if v != tc.ExpectedLarge {
				t.Fatalf("expected 'message is large' to be %v, found %v", tc.ExpectedLarge, v)
			}
		})
	}
}

func TestExtendedSQS_batchMessageIsLarge(t *testing.T) {
	cases := []struct {
		Name string
		ThresholdKb int64
		Input *sqs.SendMessageBatchInput
		ExpectedLarge bool
	}{
		{
			Name: "Batch Message Is Large",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {

		})
	}
}
