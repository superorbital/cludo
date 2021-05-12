# Secrets

resource "aws_secretsmanager_secret" "cludod_cfg" {
  name   = "${var.stack_name}-cludod-cfg"
  policy = data.aws_iam_policy_document.cludod_cfg_policy.json
}

data "aws_iam_policy_document" "cludod_cfg_policy" {
  statement {
    effect = "Allow"
    principals {
      identifiers = [var.ecs_task_execution_role_arn]
      type        = "AWS"
    }
    actions = [
      "secretsmanager:GetSecret",
      "secretsmanager:GetSecretValue"
    ]
    resources = ["*"]
  }
}

resource "aws_secretsmanager_secret_version" "cludod_cfg" {
  secret_id     = aws_secretsmanager_secret.cludod_cfg.id
  secret_string = var.cludod_cfg

  version_stages = ["AWSCURRENT"]
}

resource "aws_iam_policy" "secrets_access" {
  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetResourcePolicy",
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret",
                "secretsmanager:ListSecretVersionIds"
            ],
            "Resource": "arn:aws:secretsmanager:${var.aws_region}:${data.aws_caller_identity.current.account_id}:secret:*"
        }
    ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "secret_access" {
  role       = var.ecs_task_execution_role_name
  policy_arn = aws_iam_policy.secrets_access.arn
}


# Service definitions

resource "aws_ecs_task_definition" "cludod_server" {
  family = "cludod_server"

  requires_compatibilities = [
    "FARGATE"
  ]

  network_mode       = "awsvpc"
  execution_role_arn = aws_iam_role.cludod_role.arn

  cpu    = 256
  memory = 512

  container_definitions = <<EOF
[
  {
    "name": "cludod_server",
    "essential": true,
    "image": "superorbital/cludod:${var.cludod_version}",
    "portMappings": [{
      "hostPort": 80,
      "protocol": "tcp",
      "containerPort": 80
    }],
    "mountPoints": [],
    "volumesFrom": [],
    "cpu": 256,
    "memory": 512,
    "environment": [
      {
        "name": "PORT",
        "value": "80"
      }
    ],
    "secrets": [
      {
        "name": "CLUDOD_CONFIG",
        "valueFrom": "${aws_secretsmanager_secret.cludod_cfg.arn}"
      }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-region": "${var.aws_region}",
        "awslogs-group": "${var.stack_name}-logs",
        "awslogs-stream-prefix": "complete-ecs"
      }
    },
    "healthCheck": {
      "retries": 3,
      "command": [
        "CMD-SHELL",
        "wget -O /dev/null -o /dev/null -T 5 -t 1 http://localhost/health || exit 1"
      ],
      "timeout": 6,
      "interval": 30,
      "startPeriod": null
    }
  }
]
EOF

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_ecs_service" "cludod_server" {
  name            = "cludod_server"
  cluster         = module.ecs.this_ecs_cluster_id
  task_definition = aws_ecs_task_definition.cludod_server.arn
  # launch_type      = "FARGATE"
  platform_version = "LATEST"
  propagate_tags   = "SERVICE"

  desired_count = var.server_count

  deployment_maximum_percent         = 100
  deployment_minimum_healthy_percent = 0
  health_check_grace_period_seconds  = 10

  capacity_provider_strategy {
    base              = 0
    capacity_provider = "FARGATE_SPOT"
    weight            = 1
  }

  deployment_controller {
    type = "ECS"
  }

  load_balancer {
    target_group_arn = module.server_alb.target_group_arns[0]
    container_name   = "cludod_server"
    container_port   = 80
  }

  network_configuration {
    subnets          = module.vpc.public_subnets
    security_groups  = [aws_security_group.server_firewall.id]
    assign_public_ip = true
  }

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_security_group" "server_firewall" {
  name        = "server_firewall"
  description = "Security Group for server containers"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "HTTP from VPC"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}


# DNS/SSL

resource "aws_acm_certificate" "server" {
  domain_name       = "api.${data.aws_route53_zone.flow.name}"
  validation_method = "DNS"

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

resource "aws_route53_record" "server_validation" {
  name    = tolist(aws_acm_certificate.server.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.server.domain_validation_options)[0].resource_record_type
  zone_id = var.dns_zone_id
  records = [tolist(aws_acm_certificate.server.domain_validation_options)[0].resource_record_value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "server_cert" {
  certificate_arn         = aws_acm_certificate.server.arn
  validation_record_fqdns = [aws_route53_record.server_validation.fqdn]
}

resource "aws_route53_record" "server_alb" {
  zone_id = data.aws_route53_zone.flow.zone_id
  name    = "api"
  type    = "A"

  alias {
    name                   = module.server_alb.this_lb_dns_name
    zone_id                = module.server_alb.this_lb_zone_id
    evaluate_target_health = true
  }
}


# Load Balancers

resource "aws_security_group" "server_lb" {
  name        = "server_lb"
  description = "Allow HTTP/TLS inbound traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "TLS from internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "HTTP from internet"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description     = "Allow LoadBalancer to communicate with ECS containers"
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.server_firewall.id]
  }

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

# resource "aws_s3_bucket" "server_alb_logs" {
#   bucket = var.server_lb_log_bucket
#   acl    = "private"
# }

module "server_alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 5.0"

  name = "${var.stack_name}-server-alb"

  load_balancer_type = "application"

  vpc_id          = module.vpc.vpc_id
  subnets         = module.vpc.public_subnets
  security_groups = [aws_security_group.server_lb.id]

  # idle_timeout = 120

  # access_logs = {
  #   bucket = var.server_lb_log_bucket
  # }

  target_groups = [
    {
      name_prefix      = "srvr-"
      backend_protocol = "HTTP"
      backend_port     = 80
      target_type      = "ip"

      health_check = {
        enabled = true
        path    = "/health"
      }
    }
  ]

  https_listeners = [
    {
      port               = 443
      protocol           = "HTTPS"
      certificate_arn    = aws_acm_certificate.server.arn
      target_group_index = 0
    }
  ]

  http_tcp_listeners = [
    {
      port        = 80
      protocol    = "HTTP"
      action_type = "redirect"
      redirect = {
        port        = "443"
        protocol    = "HTTPS"
        status_code = "HTTP_301"
      }
    }
  ]

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}
