package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	awsf "github.com/edcast/one-time-secret/pkg/aws"
	"github.com/stretchr/testify/assert"
)

func TestIndexHtml(t *testing.T) {
	createEndpoint := "api/secret"
	indexHTML := Index(createEndpoint)
	assert.Contains(t, indexHTML, "var secret_text = document.getElementById(\"input_secret\").value;")              // Get secret_text var
	assert.Contains(t, indexHTML, fmt.Sprintf("xhttp.open(\"POST\", \"%s\", true);", createEndpoint))                // POST request
	assert.Contains(t, indexHTML, "xhttp.send(`{\"content\": \"${secret_text}\", \"encoding\": \"${encoding}\"}`);") // POST body
}

func TestLifeCycle(t *testing.T) {
	secretData := "1234567890123467890"
	request := CreateSecretRequest{Content: secretData}
	responseURLPrefix := "abcd"
	createResponse, createResponseErr := CreateSecret(request, responseURLPrefix)
	assert.Regexp(t, fmt.Sprintf("^%s/", responseURLPrefix), createResponse.URL)
	assert.Equal(t, http.StatusOK, createResponse.Status)
	assert.Nil(t, createResponseErr)

	secretID := strings.Replace(createResponse.URL, responseURLPrefix, "", 1)
	storedContent, getStoredContentErr := awsf.GetSecret(secretID)
	assert.Nil(t, getStoredContentErr)
	assert.NotEqual(t, secretData, storedContent)

	getResponse, getResponseErr := GetSecret(secretID)
	assert.Equal(t, http.StatusOK, getResponse.Status)
	assert.Equal(t, secretData, getResponse.Content)
	assert.Nil(t, getResponseErr)

	getResponse2, getResponseErr2 := GetSecret(secretID)
	assert.Equal(t, http.StatusInternalServerError, getResponse2.Status)
	assert.Equal(t, "", getResponse2.Content)
	assert.NotNil(t, getResponseErr2)
}

func TestBase64EncodedContent(t *testing.T) {
	secretData := `1234
1234`
	base64EncodedData := "MTIzNAoxMjM0"
	request := CreateSecretRequest{Content: base64EncodedData, Encoding: "base64"}
	responseURLPrefix := "abcd"
	createResponse, _ := CreateSecret(request, responseURLPrefix)
	secretID := strings.Replace(createResponse.URL, responseURLPrefix, "", 1)
	getResponse, _ := GetSecret(secretID)
	assert.Equal(t, secretData, getResponse.Content)
}
