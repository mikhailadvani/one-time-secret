## START Gateway creation and permission to trigger lambda

resource "aws_api_gateway_rest_api" "this" {
  count = var.lambda_enabled == true ? 1 : 0
  name  = local.project_name

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_lambda_permission" "api_gateway_trigger" {
  for_each      = local.lambda_functions
  statement_id  = "ApiGatewayInvoke${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = "${local.project_name}-${each.key}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.this[0].execution_arn}/*/*/*"
}

## END Gateway creation and permission to trigger lambda
## START Upload archive to S3 and consume in lambda

resource "aws_s3_bucket_object" "lambda_functions" {
  for_each = local.lambda_functions
  bucket   = var.bucket_name
  key      = "${local.lambda_s3_location}${local.project_name}-${each.key}"
  source   = "${var.lambda_functions_location}${local.project_name}-${each.key}.zip"
  etag     = "${filemd5("${var.lambda_functions_location}${local.project_name}-${each.key}.zip")}"
  depends_on = [aws_s3_bucket.secret_bucket]
}

resource "aws_lambda_function" "this" {
  for_each      = local.lambda_functions
  s3_bucket     = var.bucket_name
  s3_key        = aws_s3_bucket_object.lambda_functions[each.key].id
  function_name = "${local.project_name}-${each.key}"
  role          = aws_iam_role.lambda[0].arn
  handler       = "${local.project_name}-${each.key}"

  runtime = "go1.x"

  environment {
    variables = {
      BUCKET_NAME = var.bucket_name
      S3_PREFIX   = var.bucket_prefix
    }
  }
}

## END Upload archive to S3 and consume in lambda
## START IAM role & policy for operation
resource "aws_iam_role" "lambda" {
  count = var.lambda_enabled == true ? 1 : 0
  name  = var.lambda_name_prefix

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  count      = var.lambda_enabled == true ? 1 : 0
  role       = aws_iam_role.lambda[0].name
  policy_arn = aws_iam_policy.secret_object_crud.arn
}
## END IAM role & policy for operation
## START Lambda logging configurations and permissions

resource "aws_cloudwatch_log_group" "index" {
  count             = local.lambda_logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.lambda_name_prefix}-index"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "create" {
  count             = local.lambda_logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.lambda_name_prefix}-create"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "get" {
  count             = local.lambda_logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.lambda_name_prefix}-get"
  retention_in_days = 7
}

resource "aws_iam_policy" "lambda_logging" {
  count       = local.lambda_logging_enabled == true ? 1 : 0
  name        = "lambda-logging-${var.lambda_name_prefix}"
  path        = "/"
  description = "IAM policy for logging from a lambda of ${var.lambda_name_prefix}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": [
        "${aws_cloudwatch_log_group.index[0].arn}",
        "${aws_cloudwatch_log_group.create[0].arn}",
        "${aws_cloudwatch_log_group.get[0].arn}"
      ],
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logging" {
  count      = local.lambda_logging_enabled == true ? 1 : 0
  role       = "${aws_iam_role.lambda[0].name}"
  policy_arn = "${aws_iam_policy.lambda_logging[0].arn}"
}

## END Lambda logging configurations and permissions

## START Index API. No resource since it is the root_resource_id of the rest_api
resource "aws_api_gateway_method" "index" {
  count         = var.lambda_enabled == true ? 1 : 0
  rest_api_id   = "${aws_api_gateway_rest_api.this[0].id}"
  resource_id   = "${aws_api_gateway_rest_api.this[0].root_resource_id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "index" {
  count                   = var.lambda_enabled == true ? 1 : 0
  rest_api_id             = "${aws_api_gateway_rest_api.this[0].id}"
  resource_id             = "${aws_api_gateway_rest_api.this[0].root_resource_id}"
  http_method             = "${aws_api_gateway_method.index[0].http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.this["index"].arn}/invocations"
}

## END Index API

## START Create API

resource "aws_api_gateway_resource" "create" {
  count       = var.lambda_enabled == true ? 1 : 0
  rest_api_id = "${aws_api_gateway_rest_api.this[0].id}"
  parent_id   = "${aws_api_gateway_rest_api.this[0].root_resource_id}"
  path_part   = local.lambda_functions["create"]["resource"]
}

resource "aws_api_gateway_method" "create" {
  count         = var.lambda_enabled == true ? 1 : 0
  rest_api_id   = "${aws_api_gateway_rest_api.this[0].id}"
  resource_id   = "${aws_api_gateway_resource.create[0].id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "create" {
  count                   = var.lambda_enabled == true ? 1 : 0
  rest_api_id             = "${aws_api_gateway_rest_api.this[0].id}"
  resource_id             = "${aws_api_gateway_resource.create[0].id}"
  http_method             = "${aws_api_gateway_method.create[0].http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.this["create"].arn}/invocations"
}

## END Create API

## START API Deployment

resource "aws_api_gateway_deployment" "this" {
  count = var.lambda_enabled == true ? 1 : 0
  depends_on = [
    "aws_api_gateway_integration.index",
    "aws_api_gateway_integration.create",
  ]

  rest_api_id = "${aws_api_gateway_rest_api.this[0].id}"
  stage_name  = "api"
}

## END API Deployment
