resource "aws_kms_key" "this" {
  description = "KMS key used for encrypting the secret contents in S3"
}

resource "aws_kms_alias" "this" {
  name          = local.kms_key_alias
  target_key_id = aws_kms_key.this.key_id
}

data "aws_iam_role" "existing" {
  count = var.existing_iam_role == "" ? 0 : 1
  name  = var.existing_iam_role
}

resource "aws_kms_grant" "existing_role" {
  count             = var.existing_iam_role == "" ? 0 : 1
  name              = var.project_name
  key_id            = aws_kms_key.this.key_id
  grantee_principal = join("", data.aws_iam_role.existing.*.arn)
  operations        = ["Encrypt", "Decrypt"]
}

resource "aws_kms_grant" "lambda_role" {
  count             = var.lambda_enabled == true ? 1 : 0
  name              = var.project_name
  key_id            = aws_kms_key.this.key_id
  grantee_principal = aws_iam_role.lambda[0].arn
  operations        = ["Encrypt", "Decrypt"]
}

data "aws_caller_identity" "current" {}

resource "aws_kms_grant" "developer_setup" {
  count             = var.developer_setup == true ? 1 : 0
  name              = var.project_name
  key_id            = aws_kms_key.this.key_id
  grantee_principal = data.aws_caller_identity.current.arn
  operations        = ["Encrypt", "Decrypt"]
}
