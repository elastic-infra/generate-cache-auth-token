package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/elastic-infra/generate-cache-auth-token/internal/auth"
	"github.com/elastic-infra/generate-cache-auth-token/internal/config"
	"github.com/elastic-infra/generate-cache-auth-token/pkg/awsutils"
)

func main() {
	cfg := config.NewConfig()

	flag.StringVar(&cfg.UserID, "user-id", "", "IAM user ID for ElastiCache authentication (required)")
	flag.StringVar(&cfg.ReplicationGroupID, "replication-group-id", "", "ElastiCache replication group ID (required)")
	flag.StringVar(&cfg.Region, "region", cfg.Region, "AWS region (default: us-east-1)")

	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		showUsage()
		return
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		showUsage()
		os.Exit(1)
	}

	ctx := context.Background()

	awsConfig, err := awsutils.LoadAWSConfig(ctx, cfg.Region)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	tokenRequest := auth.NewIAMAuthTokenRequest(cfg.UserID, cfg.ReplicationGroupID, cfg.Region)

	token, err := tokenRequest.GenerateToken(ctx, awsConfig)
	if err != nil {
		log.Fatalf("Failed to generate IAM auth token: %v", err)
	}

	fmt.Println(token)
}

func showUsage() {
	fmt.Fprintf(os.Stderr, `ElastiCache IAM Authentication Token Generator

Usage: %s [options]

Options:
  -user-id string
        IAM user ID for ElastiCache authentication (required)
  -replication-group-id string
        ElastiCache replication group ID (required)
  -region string
        AWS region (default: ap-northeast-1)
  -help
        Show this help message

Example:
  %s -user-id iam-test-user-01 -replication-group-id iam-test-rg-01 -region us-east-1

`, os.Args[0], os.Args[0])
}
