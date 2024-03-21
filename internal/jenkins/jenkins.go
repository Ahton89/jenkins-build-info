package jenkins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sethvargo/go-githubactions"
	"io"
	"jenkins-build-info/internal/config"
	"net/http"
	"regexp"
	"time"
)

const (
	idNotFound = -1
)

func New(action *githubactions.Action, config *config.Config) Jenkins {
	return &jenkins{
		action:        action,
		configuration: config,
	}
}

func (j *jenkins) GetBuildInfo(ctx context.Context) (BuildInfo, error) {
	ticker := time.NewTicker(1 * time.Second)
	var tickerDumped bool
	var startedAt = time.Now().UTC()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("actions is cancelled: %s", ctx.Err())
		case _, ok := <-ticker.C:
			if !ok {
				return nil, errors.New("ticker channel is closed")
			}

			if time.Now().UTC().Sub(startedAt) > j.configuration.MaximumWaitTime {
				return nil, fmt.Errorf(
					"maximum wait time exceeded for the status %s",
					j.configuration.StatusName,
				)
			}

			if !tickerDumped {
				ticker.Reset(j.configuration.CheckInterval)
				tickerDumped = true
			}

			body, err := j.makeRequest(ctx)
			if err != nil {
				j.action.Warningf(
					"error making request: %s, we'll try again later after %s...",
					err,
					j.configuration.CheckInterval.Round(time.Second).String(),
				)
				continue
			}

			respObj := new(buildInfo)

			err = json.Unmarshal(body, respObj)
			if err != nil {
				j.action.Warningf(
					"error unmarshal response: %s, we'll try again later after %s...",
					err,
					j.configuration.CheckInterval.Round(time.Second).String(),
				)
				continue
			}

			id, exist := respObj.statusExist(j.configuration.StatusName)
			if exist && id != idNotFound && respObj.isSuccess(id) {
				return respObj, nil
			}

			j.action.Warningf(
				"status not found or not success, we'll try again later after %s...",
				j.configuration.CheckInterval.Round(time.Second).String(),
			)
		}
	}
}

func (j *jenkins) makeRequest(ctx context.Context) ([]byte, error) {
	var url = fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits/%s/status",
		j.configuration.GithubRepositoryOwner,
		j.configuration.GithubRepository,
		j.configuration.GithubRefName,
	)
	var client = new(http.Client)
	var resp = new(http.Response)
	var body = make([]byte, 0)

	tCtx, cancel := context.WithTimeout(ctx, j.configuration.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(tCtx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", j.configuration.InputGithubToken))
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, errors.New("status code is not 200")
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (b *buildInfo) ImageNameWithTag(config *config.Config) (string, error) {
	reBuildInfo := regexp.MustCompile(`^https://(?:.*?/+){4}.+/(\d+)/.+$`)

	id, exist := b.statusExist(config.StatusName)

	if !exist {
		return "", errors.New("status not found")
	}

	if !b.isSuccess(id) {
		return "", errors.New("status is not success")
	}

	match := reBuildInfo.FindStringSubmatch(b.Statuses[id].TargetUrl)

	if len(match) != 2 {
		return "", errors.New("build info not found or something went wrong")
	}

	return fmt.Sprintf("%s:%s-%s", config.InputImageName, config.GithubRefName, match[1]), nil
}

func (b *buildInfo) isSuccess(id int) bool {
	return b.Statuses[id].State == "success"
}

func (b *buildInfo) statusExist(statusName string) (int, bool) {
	for id, status := range b.Statuses {
		if status.Context == statusName {
			return id, true
		}
	}

	return idNotFound, false
}
