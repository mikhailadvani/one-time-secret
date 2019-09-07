package aws

import (
	"bytes"
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

// UploadSecret will upload contents to S3 and return the location
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
			}
		} else {
			log.Fatal(err.Error())
		}
		return "", err
	}
	return unprefixedKey, nil
}

// GetSecret will fetch the contents from S3
func GetSecret(unprefixedKey string) (string, error) {
	key := fmt.Sprintf("%s/%s", config.S3Prefix, unprefixedKey)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion)},
	)
	if err != nil {
		return "", err
	}
	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	contents := buf.String()
	return contents, nil
}

// DeleteSecret deletes the secret. To be called after it has been retrieved once.
func DeleteSecret(unprefixedKey string) error {
	key := fmt.Sprintf("%s/%s", config.S3Prefix, unprefixedKey)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion)},
	)
	if err != nil {
		return err
	}
	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
	}

	_, err = svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	return nil
}
