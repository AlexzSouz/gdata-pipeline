# ================================================ #
# - - - - - - - -    FUNCTIONS     - - - - - - - - #
# ================================================ #

data "aws_subnets" "subnets" {
  filter {
    name   = "vpc-id"
    values = [module.vpc.vpc_id]
  }

  depends_on = [
    module.vpc
  ]
}

data "aws_security_groups" "security_groups" {
  filter {
    name   = "vpc-id"
    values = [module.vpc.vpc_id]
  }
}

# output "vpc_subnets" {
#   description = "Subnets"
#   value       = data.aws_security_groups.security_groups

#   depends_on = [
#     module.vpc,
#     data.aws_security_groups.security_groups
#   ]
# }

resource "aws_lambda_function" "crio_functions" {
  for_each = var.crio_functions

  function_name = each.value.name
  description   = each.value.description
  role          = aws_iam_role.gdata_lambda_role.arn

  # image_uri = "${var.registry_uri}/${each.value.image_name}:${each.value.image_tag}"

  # s3_bucket = var.gdata_bucket_name # TODO : Verify Bucket
  # s3_key    = "gdata-fx.zip"

  filename = ""
  handler = "s3_sample"

  vpc_config {
    subnet_ids         = [data.aws_subnets.subnets.ids]
    security_group_ids = [data.aws_security_groups.security_groups.ids]
  }

  tags = each.value.tags

  depends_on = [
    module.vpc,
    module.s3_bucket,
    aws_iam_role.gdata_lambda_role,
    data.aws_subnets.subnets,
    data.aws_security_groups.security_groups
  ]
}
