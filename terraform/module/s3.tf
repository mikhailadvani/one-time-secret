resource "aws_s3_bucket" "secret_bucket" {
  count         = var.create_bucket == true ? 1 : 0
  bucket        = var.bucket_name
  acl           = "private"
  force_destroy = true
  versioning {
    enabled = false
  }
}

resource "aws_iam_policy" "secret_object_crud" {
  name        = "${local.project_name}-operator"
  path        = "/"
  description = "IAM policy to CRUD on ${var.bucket_name} bucket under ${var.bucket_prefix}* path for managing secrets."

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "CRUD",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObjectAcl",
                "s3:GetObject",
                "s3:DeleteObject"
            ],
            "Resource": "arn:aws:s3:::${var.bucket_name}/${var.bucket_prefix}*"
        }
    ]
}
EOF

}
