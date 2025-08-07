package config

import (
	"errors"
	"strings"
)

type Config struct {
	UserID             string
	ReplicationGroupID string
	Region             string
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.UserID) == "" {
		return errors.New("user-id cannot be empty")
	}

	if strings.TrimSpace(c.ReplicationGroupID) == "" {
		return errors.New("replication-group-id cannot be empty")
	}

	if strings.TrimSpace(c.Region) == "" {
		return errors.New("region cannot be empty")
	}

	return nil
}

func NewConfig() *Config {
	return &Config{
		Region: "ap-northeast-1",
	}
}
