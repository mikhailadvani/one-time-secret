package api

import (
	"encoding/json"
	"io"
	"log"

	"github.com/one-time-secret/pkg/aws"
)

// SecretRequest is the type for requesting secret
type SecretRequest struct {
	Content string `json:"content,omitempty"`
}

// SecretResponse is the type for requesting secret
type SecretResponse struct {
	URL      string `json:"url,omitempty"`
	ValidFor int    `json:"validFor,omitempty"`
}

// Index will return the welcome screen. Static HTML page from S3
func Index() string {
	return "Welcome!"
}

// CreateSecret will create a secret object an return a URL to access it by
func CreateSecret(requestBody io.Reader) SecretResponse {
	decoder := json.NewDecoder(requestBody)
	var request SecretRequest
	err := decoder.Decode(&request)
	if err != nil {
		log.Fatal("Unable to unmarshal request", err)
		return SecretResponse{}
	}
	secretLocation, err := aws.UploadSecret(request.Content)
	if err != nil {
		log.Fatal("Unable to upload secret", err)
		return SecretResponse{}
	}
	secretResponse := SecretResponse{URL: secretLocation}
	return secretResponse
}

// GetSecret will return the secret content stored and delete it
func GetSecret(secretID string) string {
	secretContents, err := aws.GetSecret(secretID)
	if err != nil {
		log.Fatal("Unable to get secret")
		return ""
	}
	err = aws.DeleteSecret(secretID)
	if err != nil {
		log.Fatal("Unable to delete secret")
		return ""
	}
	return secretContents
}
