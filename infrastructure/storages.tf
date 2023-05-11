# =============================================== #
# - - - - - - - -    STORAGES     - - - - - - - - #
# =============================================== #

module "s3_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket = var.gdata_bucket_name
  acl    = "private"
  # TODO : Check ACL configuration at
  #          AWS: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#permissions
  #          Terraform: https://registry.terraform.io/modules/terraform-aws-modules/s3-bucket/aws/latest#input_acl

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = true
  }

  tags = {
    IaC         = "terraform"
    Environment = var.environment
  }
}
