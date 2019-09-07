package aws

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	configf "github.com/one-time-secret/pkg/config"
)

var config = configf.LoadConfig()

// UploadSecret will upload contents to S3 and return a presigned URL
func UploadSecret(content string) (string, error) {
	unprefixedKey := uuid.New().String()
	key := fmt.Sprintf("%s/%s", config.S3Prefix, unprefixedKey)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion)},
	)
	if err != nil {
		return "", err
	}
	svc := s3.New(sess)
	input := &s3.PutObjectInput{
		ACL:    aws.String(config.S3ObjectACL),
		Body:   aws.ReadSeekCloser(strings.NewReader(content)),
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatal(aerr.Error())
				return "", aerr
			}
		} else {
			log.Fatal(err.Error())
			return "", err
		}
	}
	return unprefixedKey, nil
}
