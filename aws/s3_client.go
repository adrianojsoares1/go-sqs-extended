package big_message_aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type BigMessageS3Client interface {
	WriteBigMessage(message string) (string, error)
	ExtractBigMessage(key string) (string, error)
}

type s3Client struct {
	S3SDK *s3.S3
	Bucket string
}

func NewS3Client(s3c *s3.S3, bucket string) BigMessageS3Client {
	return &s3Client{S3SDK: s3c, Bucket: bucket}
}

func (s3c *s3Client) WriteBigMessage(message string) (string, error) {
	objectName, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("could not create unique object name for payload: %v", err)
	}
	objectNameAsString := objectName.String()
	_, err = s3c.S3SDK.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3c.Bucket),
		Key: aws.String(objectNameAsString),
		Body: strings.NewReader(message),
	})
	if err == nil {
		log.Debugf("finished uploading object %s to bucket %s", objectNameAsString, s3c.Bucket)
	}
	return objectNameAsString, err
}

func (s3c *s3Client) ExtractBigMessage(key string) (string, error) {
	output, err := s3c.S3SDK.GetObject(&s3.GetObjectInput{
		Bucket:                     aws.String(s3c.Bucket),
		Key:                        aws.String(key),
	})
	var body = new(string)
	if err != nil || output.Body == nil {
		return *body, err
	}
	sb := new(strings.Builder)
	_, err = io.Copy(sb, output.Body)
	body = aws.String(sb.String())
	return *body, err
}
