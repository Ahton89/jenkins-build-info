package config

import (
	"errors"
	"github.com/caarlos0/env"
	"os"
	"time"
)

type Config struct {
	GithubRefName         string        `env:"GITHUB_REF_NAME,required"`
	GithubRepository      string        `env:"GITHUB_REPOSITORY,required"`
	GithubRepositoryOwner string        `env:"GITHUB_REPOSITORY_OWNER,required"`
	InputGithubToken      string        `env:"INPUT_GITHUB_TOKEN"`
	InputImageName        string        `env:"INPUT_IMAGE_NAME,required"`
	RequestTimeout        time.Duration `env:"JBI_REQUEST_TIMEOUT" envDefault:"5s"`
	StatusName            string        `env:"JBI_STATUS_NAME" envDefault:"continuous-integration/jenkins/branch"`
	CheckInterval         time.Duration `env:"JBI_CHECK_INTERVAL" envDefault:"15s"`
	MaximumWaitTime       time.Duration `env:"JBI_MAXIMUM_WAIT_TIME" envDefault:"10m"`
}

func New() (Config, error) {
	configuration := Config{}
	err := env.Parse(&configuration)
	return configuration, err
}

func (c *Config) Validate() error {
	// Validating GITHUB_REPOSITORY_OWNER
	if c.GithubRepositoryOwner == "" {
		return errors.New("GITHUB_REPOSITORY_OWNER is required")
	}
	// Validating GITHUB_REPOSITORY
	if c.GithubRepository == "" {
		return errors.New("GITHUB_REPOSITORY is required")
	}
	// Validating GITHUB_REF_NAME
	if c.GithubRefName == "" {
		return errors.New("GITHUB_REF_NAME is required")
	}
	// Validating INPUT_GITHUB_TOKEN (set default value if not provided)
	if c.InputGithubToken == "" {
		defaultGithubToken := os.Getenv("GITHUB_TOKEN")
		if defaultGithubToken == "" {
			return errors.New("INPUT_GITHUB_TOKEN or GITHUB_TOKEN is required")
		}
		c.InputGithubToken = defaultGithubToken
	}
	// Validating INPUT_IMAGE_NAME
	if c.InputImageName == "" {
		return errors.New("INPUT_IMAGE_NAME is required")
	}

	return nil
}
