package go_sqs_extended

import (
	"errors"
	bma "github.com/asoares1-chwy/go-sqs-extended/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

const (
	DefaultLargeMessageSize int64 = 256 // kb
)

type ExtendedSQS struct {
	*sqs.SQS
	cfg *extendedConfigurationGhost
	s3c *s3ConfigurationGhost
}

type ExtendedConfiguration struct {
	AlwaysSendThroughS3 bool
	LargeMessageThresholdKb int64
	UseLegacyReservedAttributeName bool
	S3Configuration *S3Configuration
}

type extendedConfigurationGhost struct {
	AlwaysSendThroughS3 bool
	LargeMessageThresholdKb int64
	UseLegacyReservedAttributeName bool
}

func (ec *ExtendedConfiguration) toGhost() *extendedConfigurationGhost {
	return &extendedConfigurationGhost{
		AlwaysSendThroughS3:            ec.AlwaysSendThroughS3,
		LargeMessageThresholdKb:        ec.LargeMessageThresholdKb,
		UseLegacyReservedAttributeName: ec.UseLegacyReservedAttributeName,
	}
}

type S3Configuration struct {
	Client *s3.S3
	BucketName string
	CleanupAfterOperation bool
}

func (s3c *S3Configuration) isConfigured() bool {
	return s3c.Client != nil && s3c.BucketName != ""
}

type s3ConfigurationGhost struct {
	Configured bool
	Client bma.BigMessageS3Client
	BucketName string
	CleanupAfterOperation bool
}

func (s3c *S3Configuration) toGhost() *s3ConfigurationGhost {
	ghost := &s3ConfigurationGhost{}
	if s3c.isConfigured() {
		ghost.Configured = true
		ghost.BucketName = s3c.BucketName
		ghost.Client = bma.NewS3Client(s3c.Client, s3c.BucketName)
		ghost.CleanupAfterOperation = s3c.CleanupAfterOperation
	}
	return ghost
}

type ExtendedQueueMessage struct {
	S3BucketName string
	S3Key        string
}

func NewExtended(sqs *sqs.SQS, options *ExtendedConfiguration) (*ExtendedSQS, error) {
	if sqs == nil {
		return nil, errors.New("cannot create an extended client from a null sqs.SQS pointer")
	}
	if options == nil {
		options = &ExtendedConfiguration{S3Configuration: &S3Configuration{}}
	}
	if options.LargeMessageThresholdKb < 1 {
		log.Printf("defaulting to %dkb as large message threshold", DefaultLargeMessageSize)
		options.LargeMessageThresholdKb = DefaultLargeMessageSize
	}
	return &ExtendedSQS{
		SQS: sqs,
		cfg: options.toGhost(),
		s3c: options.S3Configuration.toGhost(),
	}, nil
}
