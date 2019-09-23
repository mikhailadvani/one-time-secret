variable region {
  description = "AWS region to create the components in."
}

variable project_name {
  default     = "one-time-secret"
  description = "Name of the complete setup. Used at several places including Lambda, KMS, etc."
}

variable bucket_name {
  description = "Name of the S3 bucket which will store the secrets."
}

variable create_bucket {
  default     = true
  description = "Whether the bucket storing the secrets is to be created or an existing bucket will be used."
}

variable bucket_prefix {
  default     = ""
  description = "The prefix under which all secrets will be created. Used to define IAM restrictions. Trailing / needed."
}

variable lambda_enabled {
  default     = false
  description = "Whether the lambda function should be deployed"
}

variable lambda_functions_location {
  description = "Location to get the zipped lambda functions from."
}

variable lambda_s3_location {
  default     = "/"
  description = "Location on S3 to upload the function archives to. Given the 4MB limit on checksum calculation by terraform. Trailing / needed."
}

variable lambda_logging_enabled {
  default     = false
  description = "Whether cloudwatch logging should be enabled or not."
}

variable existing_iam_role {
  default     = ""
  description = "Existing IAM role to be granted privilege to operate using the KMS key."
}

variable developer_setup {
  default     = false
  description = "Whether the current ARN should be given access to KMS to run tests."
}
