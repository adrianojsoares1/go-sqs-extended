package go_sqs_extended

import (
	"errors"
	bma "github.com/asoares1-chwy/go-sqs-extended/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultLargeMessageSize int64 = 262144 // bytes
	ReservedAttributeName = "ExtendedPayloadSize"
	LegacyReservedAttributeName = "SQSLargePayloadSize"
	MaximumAllowedAttributes = 9 // 10 - 1 for reserved
)

type ExtendedSQS struct {
	*sqs.SQS
	cfg *extendedConfigurationGhost
	s3c *s3ConfigurationGhost
}

type ExtendedConfiguration struct {
	AlwaysSendThroughS3            bool
	LargeMessageThreshold          int64
	UseLegacyReservedAttributeName bool
	S3Configuration                *S3Configuration
}

type extendedConfigurationGhost struct {
	AlwaysSendThroughS3            bool
	LargeMessageThreshold          int64
	UseLegacyReservedAttributeName bool
}

func (ec *ExtendedConfiguration) toGhost() *extendedConfigurationGhost {
	return &extendedConfigurationGhost{
		AlwaysSendThroughS3:            ec.AlwaysSendThroughS3,
		LargeMessageThreshold:          ec.LargeMessageThreshold,
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
		log.Debug("running with configured s3 client, messages may be sent through S3 object references")
		ghost.Configured = true
		ghost.BucketName = s3c.BucketName
		ghost.Client = bma.NewS3Client(s3c.Client, s3c.BucketName)
		ghost.CleanupAfterOperation = s3c.CleanupAfterOperation
	} else {
		log.Debug("running with unconfigured s3 client, messages can only be sent 'as is'")
	}
	return ghost
}

type extendedQueueMessage struct {
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
	if options.LargeMessageThreshold < 1 {
		log.Debugf("running with default value %d as large message threshold", DefaultLargeMessageSize)
		options.LargeMessageThreshold = DefaultLargeMessageSize
	}
	return &ExtendedSQS{
		SQS: sqs,
		cfg: options.toGhost(),
		s3c: options.S3Configuration.toGhost(),
	}, nil
}
