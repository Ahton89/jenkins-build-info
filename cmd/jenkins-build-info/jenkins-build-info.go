package main

import (
	"context"
	"github.com/sethvargo/go-githubactions"
	"jenkins-build-info/internal/config"
	"jenkins-build-info/internal/jenkins"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var buildInfo jenkins.BuildInfo

	action := githubactions.New()

	configuration, err := config.New()
	if err != nil {
		action.Fatalf("error parsing configuration: %s", err)
	}

	err = configuration.Validate()
	if err != nil {
		action.Fatalf("error validating configuration: %s", err)
	}

	jenkinsClient := jenkins.New(action, &configuration)

	buildInfo, err = jenkinsClient.GetBuildInfo(ctx)
	if err != nil {
		action.Fatalf("error getting build info: %s", err)
	}

	image, err := buildInfo.ImageNameWithTag(&configuration)
	if err != nil {
		action.Fatalf("error getting image name: %s", err)
	}

	action.SetOutput("jenkins-image-name", image)
	action.SaveState("jenkins-image-name", image)
}
