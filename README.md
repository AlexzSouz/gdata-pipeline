# GData Pipelines

This repository contains IaC to deploy AWS resources, a blank lambda function and a service to process customer orders.

**NOTE TO THE READER:** I didn't magage to create the lambda function with S3 trigger as I couldn't debug it locally, and there is a learning curve to developing AWS lambdas. I've never worked with AWS before, or with the AWS Lambda GO APIs previously so I got stuck around the tooling for local development, etc...

## Infrastructure

For more details in how to execute the IaC to provision AWS resources, please visit the [Infrastructure](./infrastructure/README.md) documentation.

**Resources Provisioned:**

- S3 Bucket
- VPC
- Policies & Rules
- Lambda Function

## Customer Saga

This service process `.csv` files related to customer orders. I've applied several coding practices and patterns and might have over-exagerated a bit.

_**NOTE:** This service can scale horizontaly by adding more instancess of it. However, some configuration and implementation code implementation is missing._

**Technologies:**

- Go (go1.19.5)
- OCI (Open Container Initiative)
  - Dockerfile (using Podman for CRI-O runtime)
  - Docker-Compose (using Podman-Compose for CRI-O runtime)
- OpenTelemetry for Metrics, Traces and Logs
- FluentTest
- FSNotify (For file watching)

**Patterns:**

- Command
- DDD Aggregate Model
- Concurrency (per customer process)

### Run

```bash
go run program.go
# Copy the .csv files to "./customer-saga/.files"
```

_**IMPORTANT:** Once the application is running, `.csv` files dropped under the folder `.files` will be processed. I ignore Orders or LineItems files, and only process everything once the Customer `.csv` file is present._
