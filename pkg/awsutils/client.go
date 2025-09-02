package awsutils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func InteractiveMFATokenProvider() (string, error) {
	fmt.Fprint(os.Stderr, "Enter MFA token: ")
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read MFA token: %w", err)
	}
	return strings.TrimSpace(token), nil
}

func LoadAWSConfig(ctx context.Context, region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = InteractiveMFATokenProvider
		}),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
