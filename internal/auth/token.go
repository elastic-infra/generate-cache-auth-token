package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

const (
	serviceName         = "elasticache"
	actionName          = "connect"
	tokenExpiryDuration = 15 * time.Minute
)

type IAMAuthTokenRequest struct {
	UserID             string
	ReplicationGroupID string
	Region             string
}

func NewIAMAuthTokenRequest(userID, replicationGroupID, region string) *IAMAuthTokenRequest {
	return &IAMAuthTokenRequest{
		UserID:             userID,
		ReplicationGroupID: replicationGroupID,
		Region:             region,
	}
}

func (r *IAMAuthTokenRequest) GenerateToken(ctx context.Context, awsConfig aws.Config) (string, error) {
	requestURL := fmt.Sprintf("http://%s/", r.ReplicationGroupID)
	u, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	query := u.Query()
	query.Set("Action", actionName)
	query.Set("User", r.UserID)
	query.Set("X-Amz-Expires", strconv.FormatInt(int64(tokenExpiryDuration/time.Second), 10))

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	creds, err := awsConfig.Credentials.Retrieve(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve AWS credentials: %w", err)
	}

	signer := v4.NewSigner()
	uri, _, err := signer.PresignHTTP(ctx, creds, req, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", serviceName, r.Region, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return strings.TrimPrefix(uri, "http://"), nil
}
