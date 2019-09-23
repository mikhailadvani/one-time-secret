package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/mikhailadvani/one-time-secret/pkg/aws"
)

// CreateSecretRequest is the type for requesting secret creation
type CreateSecretRequest struct {
	Content  string `json:"content,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

// CreateSecretResponse is the response type for requesting secret creation
type CreateSecretResponse struct {
	URL    string `json:"url,omitempty"`
	Status int    `json:"status,omitempty"`
}

// GetSecretResponse is the response type for requesting secret retrieval
type GetSecretResponse struct {
	Content string `json:"content,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// Index will return the welcome screen. Static HTML page from S3
func Index(createAPIEndpoint string) string {
	return `
	<!DOCTYPE html>
	<html>
	<body>

	<div style="width:600px; margin:0 auto;display:flex;justify-content:center;" id="container">
		<div id="create">
			<h2>Enter your secret</h2>
			<textarea id="input_secret" rows="4" cols="50"></textarea>
		  <button type="button" onclick="create()">Submit</button>
		</div>
	</div>
	<div style="width:600px; margin:0 auto;display:flex;justify-content:center;">
		<p id="generated_secret_id"></p>
	</div>

	</body>
	<script>
	  function create() {
	    var xhttp = new XMLHttpRequest();
			var secret_text = document.getElementById("input_secret").value;
			if (secret_text.indexOf("\n") == -1) {
				encoding = "utf-8"
			} else {
				secret_text = btoa(secret_text)
				encoding = "base64"
			}
	    xhttp.onreadystatechange = function() {
	      if (this.readyState == 4 && this.status == 200) {
	       document.getElementById("generated_secret_id").innerHTML = JSON.parse(this.responseText).url;
	      }
	    };
	    xhttp.open("POST", "` + createAPIEndpoint + `", true);
	    xhttp.send(` + "`" + `{"content": "${secret_text}", "encoding": "${encoding}"}` + "`" + `);
	  }
	</script>
	</html>
	`
}

// CreateSecret will create a secret object an return a URL to access it by
func CreateSecret(request CreateSecretRequest, responseURLPrefix string) (CreateSecretResponse, error) {
	var plainTextContent string
	if request.Encoding == "base64" {
		plainTextBytes, _ := base64.StdEncoding.DecodeString(request.Content)
		plainTextContent = string(plainTextBytes)
	} else {
		plainTextContent = request.Content
	}
	encryptedContent, err := aws.Encrypt(plainTextContent)
	if err != nil {
		fmt.Println(err.Error())
		return CreateSecretResponse{Status: http.StatusInternalServerError}, err
	}
	secretLocation, err := aws.UploadSecret(encryptedContent)
	if err != nil {
		return CreateSecretResponse{Status: http.StatusInternalServerError}, err
	}
	secretResponse := CreateSecretResponse{URL: fmt.Sprintf("%s/%s", responseURLPrefix, secretLocation), Status: http.StatusOK}
	return secretResponse, nil
}

// GetSecret will return the secret content stored and delete it
func GetSecret(secretID string) (GetSecretResponse, error) {
	secretContents, err := aws.GetSecret(secretID)
	if err != nil {
		return GetSecretResponse{Status: http.StatusInternalServerError}, err
	}
	decryptedContent, err := aws.Decrypt(secretContents)
	if err != nil {
		return GetSecretResponse{Status: http.StatusInternalServerError}, err
	}
	err = aws.DeleteSecret(secretID)
	if err != nil {
		return GetSecretResponse{Status: http.StatusInternalServerError}, err
	}
	return GetSecretResponse{Content: decryptedContent, Status: http.StatusOK}, nil
}
