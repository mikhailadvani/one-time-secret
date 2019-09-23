package aws

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	configf "github.com/mikhailadvani/one-time-secret/pkg/config"
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
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
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

// Encrypt will encrypt the content using KMS
func Encrypt(plainText string) (string, error) {
	svc := kms.New(session.New())
	input := &kms.EncryptInput{
		KeyId:     aws.String(config.KmsKeyAlias),
		Plaintext: []byte(plainText),
	}

	result, err := svc.Encrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				fmt.Println(kms.ErrCodeNotFoundException, aerr.Error())
			case kms.ErrCodeDisabledException:
				fmt.Println(kms.ErrCodeDisabledException, aerr.Error())
			case kms.ErrCodeKeyUnavailableException:
				fmt.Println(kms.ErrCodeKeyUnavailableException, aerr.Error())
			case kms.ErrCodeDependencyTimeoutException:
				fmt.Println(kms.ErrCodeDependencyTimeoutException, aerr.Error())
			case kms.ErrCodeInvalidKeyUsageException:
				fmt.Println(kms.ErrCodeInvalidKeyUsageException, aerr.Error())
			case kms.ErrCodeInvalidGrantTokenException:
				fmt.Println(kms.ErrCodeInvalidGrantTokenException, aerr.Error())
			case kms.ErrCodeInternalException:
				fmt.Println(kms.ErrCodeInternalException, aerr.Error())
			case kms.ErrCodeInvalidStateException:
				fmt.Println(kms.ErrCodeInvalidStateException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}
	return string(result.CiphertextBlob), nil
}

// Decrypt will encrypt the content using KMS
func Decrypt(cipherText string) (string, error) {
	svc := kms.New(session.New())
	input := &kms.DecryptInput{
		CiphertextBlob: []byte(cipherText),
	}

	result, err := svc.Decrypt(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case kms.ErrCodeNotFoundException:
				fmt.Println(kms.ErrCodeNotFoundException, aerr.Error())
			case kms.ErrCodeDisabledException:
				fmt.Println(kms.ErrCodeDisabledException, aerr.Error())
			case kms.ErrCodeInvalidCiphertextException:
				fmt.Println(kms.ErrCodeInvalidCiphertextException, aerr.Error())
			case kms.ErrCodeKeyUnavailableException:
				fmt.Println(kms.ErrCodeKeyUnavailableException, aerr.Error())
			case kms.ErrCodeDependencyTimeoutException:
				fmt.Println(kms.ErrCodeDependencyTimeoutException, aerr.Error())
			case kms.ErrCodeInvalidGrantTokenException:
				fmt.Println(kms.ErrCodeInvalidGrantTokenException, aerr.Error())
			case kms.ErrCodeInternalException:
				fmt.Println(kms.ErrCodeInternalException, aerr.Error())
			case kms.ErrCodeInvalidStateException:
				fmt.Println(kms.ErrCodeInvalidStateException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	return string(result.Plaintext), nil
}
