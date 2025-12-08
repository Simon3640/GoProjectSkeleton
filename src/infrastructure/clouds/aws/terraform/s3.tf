# S3 Bucket for Application Templates
# This bucket stores email templates and other static assets

# S3 Bucket for templates
resource "aws_s3_bucket" "templates" {
  bucket = "${local.name_prefix}-templates-${random_id.bucket_suffix.hex}"

  tags = merge(
    local.common_tags,
    {
      Name        = "${local.name_prefix}-templates"
      Description = "S3 bucket for application email templates"
    }
  )
}

# Random ID for bucket name uniqueness
resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# S3 Bucket Versioning
resource "aws_s3_bucket_versioning" "templates" {
  bucket = aws_s3_bucket.templates.id

  versioning_configuration {
    status = "Enabled"
  }
}

# S3 Bucket Server-Side Encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "templates" {
  bucket = aws_s3_bucket.templates.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# S3 Bucket Public Access Block (keep bucket private)
resource "aws_s3_bucket_public_access_block" "templates" {
  bucket = aws_s3_bucket.templates.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 Bucket Lifecycle Configuration
resource "aws_s3_bucket_lifecycle_configuration" "templates" {
  bucket = aws_s3_bucket.templates.id

  rule {
    id     = "delete-old-versions"
    status = "Enabled"

    filter {
      prefix = ""
    }

    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}

# Upload email templates to S3
# This synchronizes templates from the codebase to S3
resource "aws_s3_object" "email_templates" {
  for_each = fileset("${path.module}/../../../../application/shared/templates/emails", "*.gohtml")

  bucket = aws_s3_bucket.templates.id
  key    = "templates/emails/${each.value}"
  source = "${path.module}/../../../../application/shared/templates/emails/${each.value}"

  etag = filemd5("${path.module}/../../../../application/shared/templates/emails/${each.value}")

  content_type = "text/html"

  tags = merge(
    local.common_tags,
    {
      Name = "email-template-${each.value}"
      Type = "email-template"
    }
  )
}

# IAM Policy for Lambda functions to read from S3 templates bucket
resource "aws_iam_policy" "s3_templates_read" {
  name        = "${local.name_prefix}-s3-templates-read"
  description = "Policy to allow Lambda functions to read templates from S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:ListBucket"
        ]
        Resource = [
          aws_s3_bucket.templates.arn,
          "${aws_s3_bucket.templates.arn}/*"
        ]
      }
    ]
  })

  tags = local.common_tags
}
