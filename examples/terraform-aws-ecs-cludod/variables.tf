variable "aws_region" {
  type        = string
  default     = "us-west-1"
  description = "AWS region to deploy cludod to"
}

variable "stack_name" {
  type        = string
  default     = "cludod"
  description = "Name prefix for created resources"
}

variable "vpc_cidr" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR for cludod VPC"
}

variable "vpc_azs" {
  type        = list(string)
  default     = ["us-west-1b", "us-west-1c"]
  description = "List of AWS availability zones to deploy cludod VPC to"
}

variable "vpc_public_subnets" {
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24"]
  description = "List of subnet CIDRs to deploy cludod to. Needs to match the size of vpc_azs"
}

variable "cludod_version" {
  type        = string
  default     = "latest"
  description = "Version of the superorbital/cludod container to deploy"
}

variable "dns_zone_id" {
  type        = string
  description = "Route53 zone id of the deployed application"
}

variable "server_count" {
  type        = number
  default     = 1
  description = "Replica count for the cludod server"
}

variable "cludod_cfg" {
  type        = string
  description = "Cludod YAML configuration as a string. Will be stored in an AWS secrets manager secret"
}
