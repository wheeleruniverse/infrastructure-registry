
resource "aws_apigatewayv2_api" "this" {
  /*
   * Note:
   *
   * If the body argument is provided, the OpenAPI specification will be used to configure the integrations and route
   * for the HTTP API. If this argument is provided, the following resources should not be managed as separate ones,
   * as updates may cause manual resource updates to be overwritten:
   * - aws_apigatewayv2_integration
   * - aws_apigatewayv2_route
   *
   * https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_api
   */
  body          = file("../../swagger/api.json")
  name          = var.domain
  protocol_type = "HTTP"
  tags          = local.tags
}

resource "aws_apigatewayv2_deployment" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  description = var.domain
  triggers = {
    redeployment = sha1(jsonencode(aws_apigatewayv2_api.this.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_stage" "this" {
  api_id        = aws_apigatewayv2_api.this.id
  deployment_id = aws_apigatewayv2_deployment.this.id
  name          = var.environment
  tags          = local.tags
}

resource "aws_cloudfront_distribution" "this" {
  aliases             = local.aws_fqdn_list
  enabled             = true
  comment             = "${var.domain} distribution"
  default_root_object = "index.html"
  is_ipv6_enabled     = true
  price_class         = "PriceClass_100"
  tags                = local.tags

  default_cache_behavior {
    allowed_methods = [
      "DELETE",
      "GET",
      "HEAD",
      "OPTIONS",
      "PATCH",
      "POST",
      "PUT"
    ]
    cached_methods = [
      "GET",
      "HEAD"
    ]
    compress               = true
    target_origin_id       = aws_s3_bucket.web.id
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }
  origin {
    domain_name = aws_s3_bucket.web.bucket_regional_domain_name
    origin_id   = aws_s3_bucket.web.id

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.this.cloudfront_access_identity_path
    }
  }
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
  viewer_certificate {
    acm_certificate_arn            = "arn:aws:acm:${var.region}:${var.account}:certificate/5bed8c4a-fb53-4a2e-800d-64740eaf6c7b"
    cloudfront_default_certificate = false
    minimum_protocol_version       = var.tls_version
    ssl_support_method             = "sni-only"
  }
}

resource "aws_cloudfront_origin_access_identity" "this" {
  comment = "${var.domain} oai"
}

resource "aws_codeartifact_domain" "this" {
  domain = var.domain
  tags   = local.tags
}

resource "aws_codeartifact_repository" "this" {
  domain     = aws_codeartifact_domain.this.domain
  repository = var.project
  tags       = local.tags
}

resource "aws_codeartifact_repository_permissions_policy" "this" {
  domain          = aws_codeartifact_domain.this.domain
  repository      = aws_codeartifact_repository.this.repository
  policy_document = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "codeartifact:DescribePackageVersion",
                "codeartifact:DescribeRepository",
                "codeartifact:GetPackageVersionReadme",
                "codeartifact:GetRepositoryEndpoint",
                "codeartifact:ListPackageVersionAssets",
                "codeartifact:ListPackageVersionDependencies",
                "codeartifact:ListPackageVersions",
                "codeartifact:ListPackages",
                "codeartifact:ReadFromRepository"
            ],
            "Effect": "Allow",
            "Principal": "*",
            "Resource": "${aws_codeartifact_repository.this.arn}"
        }
    ]
}
EOF
}

resource "aws_dynamodb_table" "this" {
  billing_mode   = "PROVISIONED"
  count          = length(var.tables)
  hash_key       = "id"
  name           = var.tables[count.index]
  read_capacity  = 1
  tags           = local.tags
  write_capacity = 1

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_ecr_repository" "this" {
  image_tag_mutability = "IMMUTABLE"
  name                 = var.domain
  tags                 = local.tags

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_lambda_function" "this" {
  function_name = var.domain
  image_uri     = "${var.account}.dkr.ecr.${var.region}.amazonaws.com/${var.domain}:${var.image}"
  memory_size   = 512
  package_type  = "Image"
  role          = var.lambda_role
  tags          = local.tags
  timeout       = 30
}

resource "aws_s3_bucket" "app" {
  acl    = "private"
  bucket = "${var.domain}-app"
  tags   = local.tags
}

resource "aws_s3_bucket" "iac" {
  acl    = "public-read"
  bucket = "${var.domain}-iac"
  tags   = local.tags
}

resource "aws_s3_bucket" "web" {
  acl    = "private"
  bucket = "${var.domain}-web"
  tags   = local.tags
}
