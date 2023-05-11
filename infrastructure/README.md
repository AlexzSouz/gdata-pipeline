# GData Pipeline Infrastructre

IaC using Terraform to deploy AWS resources.

## Initial Configuration

**CLI Setup:**

Check the [cli@v2] documentation for CLI reference, and for more details in how to sign-in to your AWS CLI, please visit the [sso-configure-profile-token-auto-sso] documentation.

```bash
# https://eu-central-1.console.aws.amazon.com/singlesignon/identity/home?region=eu-central-1#!/instances/6987debe3b37bbf5/dashboard
aws configure sso --no-browser
```

## IAM Service Accounts Setup

Check the [AWS Policy List] for details when configuring AWS Group Policies.

**Creating IaC Provisioning SA, Groups and Keys:**

```bash
aws iam create-group \
  --group-name <group-name>

aws iam attach-group-policy \
  --group-name <group-name> \
  --policy-arn <policy-arn>

aws iam create-user \
  --user-name <user-name> \
  --tags Key=type,Value=iac-provisioner \
  --region <region-name>

aws iam add-user-to-group \
  --user-name <user-name> \
  --group-name <group-name>

# SAMPLE
aws iam create-group \
  --group-name iac-maintainers

policies=(\
  arn:aws:iam::aws:policy/AdministratorAccess \
  arn:aws:iam::aws:policy/IAMFullAccess \
  arn:aws:iam::aws:policy/AmazonS3FullAccess \
  arn:aws:iam::aws:policy/AWSLambda_FullAccess \
  arn:aws:iam::aws:policy/job-function/NetworkAdministrator)
for policy in "${policies[@]}"; do
  aws iam attach-group-policy \
    --group-name iac-maintainers \
    --policy-arn "$policy"
done

aws iam create-user \
  --user-name iac-tf-provisioner \
  --tags Key=type,Value=iac-provisioner \
  --region eu-central-1

aws iam add-user-to-group \
  --user-name iac-tf-provisioner \
  --group-name iac-maintainers

aws iam create-access-key \
  --user-name iac-tf-provisioner
```

**AWS Environment Variables for Terraform CLI:**

```bash
export AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID"
export AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY"
```

## TF Execution

**Init & Plan Terraform:**

```bash
# Enables Debugging
export TF_LOG=JSON

# Initializing configuration
terraform init -backend-config="bucket=<aws-bucket-name>"

# Planning configuration
terraform plan -var-file ./values/<environment>.tfvars -out <identifier>-<environment>.tfplan
```

**Apply Terraform:**

```bash
terraform apply -auto-approve <identifier>-<environment>.tfplan
```

[//]: ----------------------------------

[sso-configure-profile-token-auto-sso]: https://docs.aws.amazon.com/cli/latest/userguide/sso-configure-profile-token.html#sso-configure-profile-token-auto-sso
[cli@v2]: https://awscli.amazonaws.com/v2/documentation/api/latest/reference/index.html
[AWS Policy List]: https://docs.aws.amazon.com/aws-managed-policy/latest/reference/policy-list.html
