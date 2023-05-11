# ====================================================== #
# - - - - - - - -    VARIABLES SETUP     - - - - - - - - #
# ====================================================== #

#
# General
#
variable "environment" {
  description = "Operating environment"
  type        = string
  default     = "development"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "registry_uri" {
  description = "OCI container registry URI"
  type        = string
  default     = "quay.io"
}

#
# Virtual Private Cloud (VPC)
#
variable "vpc_name" {
  description = "VPC name"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR"
  type        = string
  default     = "10.0.0.0/16"
}

variable "vpc_subnet_azs" {
  description = "VPC subnet regions"
  type        = list(string)
}

#
# S3 Bucket
#
variable "gdata_bucket_name" {
  description = "S3 data pipeline bucket name"
  type        = string
}

#
# Functions
#
variable "crio_functions" {
  description = "Lambda CRI-O functions"
  type = list(object({
    name        = string
    description = string
    image_name  = string
    image_tag   = string
    tags        = map(string)
  }))
}
