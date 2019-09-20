module "one_time_secret" {
  source                    = "../module"
  region                    = "eu-west-1"
  bucket_name               = "one-time-secret-secrets"
  create_bucket             = true
  bucket_prefix             = "data/"
  lambda_enabled            = true
  lambda_functions_location = "${path.module}/../../build/"
  lambda_logging_enabled    = true
}
