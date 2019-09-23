output bucket_operate_policy {
  value = module.one_time_secret.bucket_operate_policy
}

output index_endpoint {
  value = module.one_time_secret.index_endpoint
}

output create_endpoint {
  value = module.one_time_secret.create_endpoint
}

output get_endpoint {
  value = module.one_time_secret.get_endpoint
}

output kms_config {
  value = module.one_time_secret.kms_config
}
