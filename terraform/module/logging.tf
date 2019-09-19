resource "aws_cloudwatch_log_group" "index" {
  count = var.logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.stack_name}-index"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "create" {
  count = var.logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.stack_name}-create"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "get" {
  count = var.logging_enabled == true ? 1 : 0
  name              = "/aws/lambda/${var.stack_name}-get"
  retention_in_days = 7
}

resource "aws_iam_policy" "lambda_logging" {
  count = var.logging_enabled == true ? 1 : 0
  name = "lambda-logging-${var.stack_name}"
  path = "/"
  description = "IAM policy for logging from a lambda of ${var.stack_name}"

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
        ${aws_cloudwatch_log_group.index.arn},
        ${aws_cloudwatch_log_group.create.arn},
        ${aws_cloudwatch_log_group.get.arn},
      ],
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logging" {
  count = var.logging_enabled == true ? 1 : 0
  role = "${aws_iam_role.lambda.name}"
  policy_arn = "${aws_iam_policy.lambda_logging.arn}"
}
