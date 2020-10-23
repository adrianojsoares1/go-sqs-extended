package big_message_aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"log"
	"strings"
)

type BigMessageS3Client interface {
	WriteBigMessage(message string) (string, error)
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
	log.Printf("finished uploading object %s to bucket %s", objectNameAsString, s3c.Bucket)
	return objectNameAsString, err
}
