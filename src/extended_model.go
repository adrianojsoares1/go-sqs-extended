package go_sqs_extended

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

const (
	DefaultLargeMessageSize = 256 // kb
)

type ExtendedSQS struct {
	*sqs.SQS
	Configuration *ExtendedConfiguration
}

type ExtendedConfiguration struct {
	AlwaysSendThroughS3 bool
	LargeMessageThresholdKb int64
	UseLegacyReservedAttributeName bool
	S3Configuration *S3Configuration
}

type S3Configuration struct {
	Client *s3.S3
	BucketName string
	CleanupAfterOperation bool
}

type ExtendedQueueMessage struct {
	S3BucketName string
	S3Key string
}

func (s3c *S3Configuration) isConfigured() bool {
	return s3c.Client != nil && s3c.BucketName != ""
}

func NewExtended(sqs *sqs.SQS, options *ExtendedConfiguration) (*ExtendedSQS, error) {
	if sqs == nil {
		return nil, errors.New("cannot create an extended client from a null sqs.SQS pointer")
	}
	if options == nil {
		options = &ExtendedConfiguration{}
	}
	if options.LargeMessageThresholdKb < 1 {
		log.Printf("defaulting to %dkb as large message threshold", DefaultLargeMessageSize)
		options.LargeMessageThresholdKb = DefaultLargeMessageSize
	}
	return &ExtendedSQS{
		SQS:             sqs,
		Configuration:   options,
	}, nil
}
