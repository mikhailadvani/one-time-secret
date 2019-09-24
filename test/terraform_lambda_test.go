package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/aws"
)

func TestLambdaCreateAndGetApis(t *testing.T) {
	t.Parallel()
	bucketName := fmt.Sprintf("terratest-one-time-secret-%s", strings.ToLower(random.UniqueId()))
	projectName := fmt.Sprintf("terratest-one-time-secret-%s", strings.ToLower(random.UniqueId()))

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform/module",
		VarFiles:     []string{"../../test/test.tfvars"},
		Vars: map[string]interface{}{
			"region":       awsRegion,
			"bucket_name":  bucketName,
			"project_name": projectName,
		},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	createEndpointSpec := terraform.OutputMap(t, terraformOptions, "create_endpoint")
	assert.Equal(t, "POST", createEndpointSpec["method"])
	assert.Regexp(t, "/api/secret$", createEndpointSpec["url"])

	getEndpointSpec := terraform.OutputMap(t, terraformOptions, "get_endpoint")
	assert.Equal(t, "GET", getEndpointSpec["method"])
	assert.Regexp(t, fmt.Sprintf("%s$", regexp.QuoteMeta(`/api/secret/{secretID+}`)), getEndpointSpec["url"])

	indexEndpointSpec := terraform.OutputMap(t, terraformOptions, "index_endpoint")
	assert.Equal(t, "GET", indexEndpointSpec["method"])
	assert.Regexp(t, "/api$", indexEndpointSpec["url"])

	kmsConfig := terraform.OutputMap(t, terraformOptions, "kms_config")
	assert.Equal(t, fmt.Sprintf("alias/%s", projectName), kmsConfig["alias"])

	bucketOperatePolicy := terraform.Output(t, terraformOptions, "bucket_operate_policy")
	assert.Regexp(t, fmt.Sprintf("/%s-operator", projectName), bucketOperatePolicy)

	indexResponseBody := httpGetWithNoErrors(t, indexEndpointSpec["url"], 200)
	assert.Regexp(t, "Enter your secret", indexResponseBody)

	secretContent := "My very secret text"
	var jsonStr = fmt.Sprintf(`{"content": "%s", "encoding": "utf-8"}`, secretContent)
	postRequestBody := strings.NewReader(jsonStr)

	postResponseBody := httpPostWithNoErrors(t, createEndpointSpec["url"], "application/json", postRequestBody, 200)
	getSecretURL := getStringValueFromJSONResponseWithNoErrors(t, postResponseBody, "url")

	getResponseBody := httpGetWithNoErrors(t, getSecretURL, 200)
	assert.Equal(t, secretContent, string(getResponseBody))
}

func getStringValueFromJSONResponseWithNoErrors(t *testing.T, body []byte, key string) string {
	var responseJSON map[string]interface{}
	err := json.Unmarshal(body, &responseJSON)
	assert.Nil(t, err)
	i := responseJSON[key]
	return interfaceToStringWithNoErrors(t, i)
}

func httpPostWithNoErrors(t *testing.T, url string, contentType string, body io.Reader, expectedResponseCode int) []byte {
	response, err := http.Post(url, contentType, body)
	assert.Nil(t, err)
	defer response.Body.Close()
	assert.Equal(t, expectedResponseCode, response.StatusCode)
	responseBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	return responseBody
}

func httpGetWithNoErrors(t *testing.T, url string, expectedResponseCode int) string {
	response, err := http.Get(string(url))
	assert.Nil(t, err)
	defer response.Body.Close()
	assert.Equal(t, expectedResponseCode, response.StatusCode)
	responseBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	return string(responseBody)
}

func interfaceToStringWithNoErrors(t *testing.T, i interface{}) string {
	s, ok := i.(string)
	assert.True(t, ok)
	return s
}
