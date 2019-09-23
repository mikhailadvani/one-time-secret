package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestLifeCycle(t *testing.T) {
	inputText := "some-random-string-to-be-uploaded"
	s3Key, _ := UploadSecret(inputText)
	assert.NotEmpty(t, s3Key)

	storedText, _ := GetSecret(s3Key)
	assert.Equal(t, inputText, storedText)

	DeleteSecret(s3Key)

	deletedText, err := GetSecret(s3Key)
	aerr, _ := err.(awserr.Error)

	assert.Equal(t, "", deletedText)
	assert.Equal(t, err.(awserr.Error), err)
	assert.Equal(t, s3.ErrCodeNoSuchKey, aerr.Code())
}

func TestEncryption(t *testing.T) {
	inputText := "some-random-text-to-be-encrypted"
	encryptedText, encryptionError := Encrypt(inputText)
	assert.Nil(t, encryptionError)
	assert.NotEqual(t, inputText, encryptedText)

	plainText, decryptionError := Decrypt(encryptedText)
	assert.Nil(t, decryptionError)
	assert.Equal(t, inputText, plainText)
}
