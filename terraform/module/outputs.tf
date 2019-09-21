output bucket_operate_policy {
  value       = aws_iam_policy.secret_object_crud.arn
  description = "Policy ARN of the IAM policy allowing operations on the bucket objects."
}

output create_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}/${local.lambda_functions["create"]["resource"]}"
    method = aws_api_gateway_method.create[0].http_method
  }
}

output get_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}/${local.lambda_functions["get"]["resource"]}"
    method = aws_api_gateway_method.get[0].http_method
  }
}

output index_endpoint {
  value = {
    url    = "${aws_api_gateway_deployment.this[0].invoke_url}/${local.lambda_functions["index"]["resource"]}"
    method = aws_api_gateway_method.index[0].http_method
  }
}
