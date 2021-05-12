data "aws_route53_zone" "cludod" {
  zone_id = var.dns_zone_id
}
