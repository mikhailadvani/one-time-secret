locals {
  lambda_logging_enabled = var.lambda_enabled && var.lambda_logging_enabled
  lambda_functions       = var.lambda_enabled == true ? { "index" : { "resource" : "", "method" : "get" }, "create" : { "resource" : "secret", "method" : "post" }, "get" : { "resource" : "{secretID+}", "method" : "get" } } : {}
  lambda_s3_location     = var.lambda_s3_location == "/" ? "" : var.lambda_s3_location
  kms_key_alias          = "alias/${var.project_name}"
  project_name           = "one-time-secret"
}
