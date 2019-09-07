package config

import (
	"fmt"
	"os"
)

// Config holds all configuration params that will be read either from config.yaml or environment variables
type Config struct {
	AwsRegion   string `json:"awsRegion,omitempty"`
	BucketName  string `json:"bucketName,omitempty"`
	S3Prefix    string `json:"s3Prefix,omitempty"`
	S3ObjectACL string `json:"s3ObjectACL,omitempty"`
}

// LoadConfig loads config attributes from env variables and returns Config object merged with defaults
func LoadConfig() Config {
	awsRegion := getMandatoryEnvironmentVariable("AWS_REGION")
	bucketName := getMandatoryEnvironmentVariable("BUCKET_NAME")
	s3Prefix := getOptionalEnvironmentVariable("S3_PREFIX", "")
	s3objectACL := getOptionalEnvironmentVariable("S3_OBJECT_ACL", "bucket-owner-full-control")
	config := Config{
		AwsRegion:   awsRegion,
		BucketName:  bucketName,
		S3Prefix:    s3Prefix,
		S3ObjectACL: s3objectACL,
	}
	return config
}

func getMandatoryEnvironmentVariable(key string) string {
	value, err := getMandatoryEnvironmentVariableE(key)
	if err != nil {
		panic(err.Error())
	}
	return value
}

func getMandatoryEnvironmentVariableE(key string) (string, error) {
	value, defined := os.LookupEnv(key)
	if !defined {
		return "", fmt.Errorf("%s environment variable not defined", key)
	}
	return value, nil
}

func getOptionalEnvironmentVariable(key string, fallback string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return fallback
	}
	return value
}
