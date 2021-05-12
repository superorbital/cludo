provider "aws" {
  region = var.aws_region
}


# VPC

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.stack_name}-vpc"
  cidr = var.vpc_cidr

  azs            = var.vpc_azs
  public_subnets = var.vpc_public_subnets

  enable_nat_gateway   = true
  enable_vpn_gateway   = true
  enable_dns_hostnames = true

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}


# ECS/Fargate Cluster

module "ecs" {
  source = "terraform-aws-modules/ecs/aws"

  name = "${var.stack_name}-ecs"

  container_insights = true

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy = [
    {
      capacity_provider = "FARGATE_SPOT"
      weight            = 1
    }
  ]

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}


# IAM

resource "aws_cloudwatch_log_group" "cludod" {
  name              = "${var.stack_name}-logs"
  retention_in_days = 1

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_iam_role" "cludod_role" {
  name = "${var.stack_name}-role"

  # May be necessary: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role
  # force_detach_policies = true

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_iam_policy" "cludod_policy" {
  name        = "${var.stack_name}-policy"
  description = "A policy for running cludod ECS containers on Fargate"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cludod_policy_attachment" {
  role       = aws_iam_role.cludod_role.name
  policy_arn = aws_iam_policy.cludod_policy.arn
}
