# =============================================== #
# - - - - - -    POLICIES & ROLES     - - - - - - #
# =============================================== #

data "aws_iam_policy_document" "gdata_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = [
      "sts:AssumeRole",
    ]
  }
}

resource "aws_iam_role" "gdata_lambda_role" {
  name               = "gdata-pipeline-role"
  assume_role_policy = data.aws_iam_policy_document.gdata_assume_role.json
}

resource "aws_iam_role_policy_attachment" "gdata_lambda_basicexecution_policy" {
  role       = aws_iam_role.gdata_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "gdata_lambda_s3execution_policy" {
  role       = aws_iam_role.gdata_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonS3ObjectLambdaExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "AWSLambdaVPCAccessExecutionRole" {
  role       = aws_iam_role.gdata_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

#
# NOTES
#
#   Actions docs: https://docs.aws.amazon.com/AmazonS3/latest/API/API_Operations.html
#     s3:GetObject: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
