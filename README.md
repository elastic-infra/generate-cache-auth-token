# ElastiCache IAM Authentication Token Generator

A Golang implementation of AWS ElastiCache IAM authentication token generator, similar to the Java-based [elasticache-iam-auth-demo-app](https://github.com/aws-samples/elasticache-iam-auth-demo-app).

## Overview

This tool generates IAM authentication tokens for ElastiCache for Redis using AWS SigV4 signing. The generated tokens can be used as passwords for IAM-based authentication to ElastiCache clusters.

## Requirements

- Go 1.21 or later
- AWS credentials configured (via AWS CLI, environment variables, or IAM roles)
- ElastiCache for Redis version 7.0 or higher with TLS enabled
- The application must run in the same VPC as the ElastiCache cluster

## Installation

```bash
git clone https://github.com/elastic-infra/generate-cache-auth-token.git
cd generate-cache-auth-token
go build -o elasticache-token ./cmd/elasticache-token
```

## Usage

### Basic Usage

```bash
./elasticache-token -user-id <user-id> -replication-group-id <replication-group-id> -region <region>
```

### Example

```bash
./elasticache-token -user-id iam-test-user-01 -replication-group-id iam-test-rg-01 -region us-east-1
```

### Command Line Options

- `-user-id`: IAM user ID for ElastiCache authentication (required)
- `-replication-group-id`: ElastiCache replication group ID (required)
- `-region`: AWS region (default: ap-northeast-1)
- `-help`: Show help message

## Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
generate-cache-auth-token/
├── cmd/
│   └── elasticache-token/
│       └── main.go              # Main application entry point
├── internal/
│   ├── auth/
│   │   └── token.go             # IAM authentication token generation logic
│   └── config/
│       └── config.go            # Configuration structure and validation
├── pkg/
│   └── awsutils/
│       └── client.go            # AWS SDK utilities
├── go.mod                       # Go module definition
├── go.sum                       # Dependency lock file
└── README.md                    # This file
```

## How It Works

1. The tool uses AWS SDK for Go v2 to obtain AWS credentials
2. Creates an HTTP GET request to the ElastiCache service endpoint
3. Signs the request using AWS SigV4 signature
4. Returns the signed URL (without http:// prefix) as the authentication token
5. The token is valid for 15 minutes

## AWS Credentials

The tool uses the AWS SDK's default credential provider chain, which looks for credentials in this order:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. AWS credentials file (`~/.aws/credentials`)
3. IAM roles for EC2 instances
4. IAM roles for tasks (ECS/Fargate)

## License

This project is released under the MIT License.
