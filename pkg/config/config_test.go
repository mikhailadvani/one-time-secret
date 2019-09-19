package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigSimple(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("BUCKET_NAME")
	os.Unsetenv("S3_PREFIX")
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("BUCKET_NAME", "some-random-bucket")
	defer os.Unsetenv("AWS_REGION")
	defer os.Unsetenv("BUCKET_NAME")
	expectedConfig := Config{
		AwsRegion:   "eu-west-1",
		BucketName:  "some-random-bucket",
		S3Prefix:    "",
		S3ObjectACL: "bucket-owner-full-control",
		BaseURL:     "http://localhost:8080",
	}
	actualConfig := LoadConfig()
	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf(`LoadConfig() got:
				%v
				want:
				%v`, actualConfig, expectedConfig)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("BUCKET_NAME", "some-random-bucket2")
	os.Setenv("S3_PREFIX", "secrets")
	os.Setenv("BASE_URL", "https://api-gateway-id.execute-api.eu-west-1.amazonaws.com/stage")

	defer os.Unsetenv("AWS_REGION")
	defer os.Unsetenv("BUCKET_NAME")
	defer os.Unsetenv("S3_PREFIX")
	defer os.Unsetenv("BASE_URL")

	expectedConfig := Config{
		AwsRegion:   "eu-west-1",
		BucketName:  "some-random-bucket2",
		S3Prefix:    "secrets",
		S3ObjectACL: "bucket-owner-full-control",
		BaseURL:     "https://api-gateway-id.execute-api.eu-west-1.amazonaws.com/stage",
	}
	actualConfig := LoadConfig()
	if !reflect.DeepEqual(expectedConfig, actualConfig) {
		t.Errorf(`LoadConfig() got:
				%v
				want:
				%v`, actualConfig, expectedConfig)
	}
}

func TestPanicOnBucketNameNotDefined(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-west-1")
	defer os.Unsetenv("AWS_REGION")
	os.Unsetenv("BUCKET_NAME")
	assert.PanicsWithValue(t, "BUCKET_NAME environment variable not defined", func() { LoadConfig() })
}

func TestPanicOnAwsRegionNotDefined(t *testing.T) {
	os.Setenv("BUCKET_NAME", "some-random-bucket")
	defer os.Unsetenv("BUCKET_NAME")
	os.Unsetenv("AWS_REGION")
	assert.PanicsWithValue(t, "AWS_REGION environment variable not defined", func() { LoadConfig() })
}
