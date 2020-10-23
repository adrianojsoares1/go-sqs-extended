package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"testing"
)

func TestExtendedSQS_messageIsLarge(t *testing.T) {
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
