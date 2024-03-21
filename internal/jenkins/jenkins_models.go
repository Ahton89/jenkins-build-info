package jenkins

import (
	"context"
	"github.com/sethvargo/go-githubactions"
	"jenkins-build-info/internal/config"
)

type jenkins struct {
	configuration *config.Config
	action        *githubactions.Action
}

type Jenkins interface {
	GetBuildInfo(ctx context.Context) (BuildInfo, error)
}

type buildInfo struct {
	State    string `json:"state"`
	Statuses []struct {
		State     string `json:"state"`
		Context   string `json:"context"`
		TargetUrl string `json:"target_url"`
	} `json:"statuses"`
}

type BuildInfo interface {
	ImageNameWithTag(config *config.Config) (string, error)
}
