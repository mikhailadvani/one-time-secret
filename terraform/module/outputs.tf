output bucket_operate_policy {
  value       = aws_iam_policy.secret_object_crud.arn
  description = "Policy ARN of the IAM policy allowing operations on the bucket objects."
}

output create_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}${aws_api_gateway_resource.create[0].path}"
    method = aws_api_gateway_method.create[0].http_method
  }
}

output get_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}${aws_api_gateway_resource.get[0].path}"
    method = aws_api_gateway_method.get[0].http_method
  }
}

output index_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}"
    method = aws_api_gateway_method.index[0].http_method
  }
}

output kms_config {
  value = {
    alias = local.kms_key_alias
    grantees = compact([
      aws_kms_grant.existing_role[0].grantee_principal,
      aws_kms_grant.lambda_role[0].grantee_principal,
      aws_kms_grant.developer_setup[0].grantee_principal,
      ],
    ),
  }
}
