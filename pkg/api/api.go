package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mikhailadvani/one-time-secret/pkg/aws"
)

// CreateSecretRequest is the type for requesting secret creation
type CreateSecretRequest struct {
	Content string `json:"content,omitempty"`
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
func Index() string {
	return `
	<!DOCTYPE html>
	<html>
	<body>

	<div style="width:400px; margin:0 auto;display:flex;justify-content:center;" id="container">
		<div id="create">
			<h2>Enter your secret</h2>
		  <input id="input_secret"></input>
		  <button type="button" onclick="create()">Submit</button>
		</div>
		<div>
			<p id="generated_secret_id"></p>
		</div>
	</div>
	</body>
	<script>
	  function create() {
	    var xhttp = new XMLHttpRequest();
			var secret_text = document.getElementById("input_secret").value;
	    xhttp.onreadystatechange = function() {
	      if (this.readyState == 4 && this.status == 200) {
	       document.getElementById("generated_secret_id").innerHTML = JSON.parse(this.responseText).url;
	      }
	    };
	    xhttp.open("POST", "/api/v1/secret", true);
	    xhttp.send(` + "`" + `{"content": "${secret_text}"}` + "`" + `);
	  }
	</script>
	</html>
	`
}

// CreateSecret will create a secret object an return a URL to access it by
func CreateSecret(requestBody io.Reader, prefix string) (CreateSecretResponse, error) {
	decoder := json.NewDecoder(requestBody)
	var request CreateSecretRequest
	err := decoder.Decode(&request)
	if err != nil {
		return CreateSecretResponse{Status: http.StatusInternalServerError}, err
	}
	secretLocation, err := aws.UploadSecret(request.Content)
	if err != nil {
		return CreateSecretResponse{Status: http.StatusInternalServerError}, err
	}
	secretResponse := CreateSecretResponse{URL: fmt.Sprintf("%s/%s", prefix, secretLocation), Status: http.StatusOK}
	return secretResponse, nil
}

// GetSecret will return the secret content stored and delete it
func GetSecret(secretID string) (GetSecretResponse, error) {
	secretContents, err := aws.GetSecret(secretID)
	if err != nil {
		return GetSecretResponse{Status: http.StatusInternalServerError}, err
	}
	err = aws.DeleteSecret(secretID)
	if err != nil {
		return GetSecretResponse{Status: http.StatusInternalServerError}, err
	}
	return GetSecretResponse{Content: secretContents, Status: http.StatusOK}, nil
}
