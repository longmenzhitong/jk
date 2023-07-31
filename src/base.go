package jk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LastBuild struct {
	Stages []struct {
		Name     string `json:"name"`
		Status   string `json:"status"`
		Complete bool   `json:"complete"`
	} `json:"stages"`
}

func getLastBuild(job string, successful bool) (*LastBuild, error) {
	var statusUrl string
	if successful {
		statusUrl = fmt.Sprintf("%s/job/%s/lastSuccessfulBuild/wfapi/describe", Config.Jenkins.Url, job)
	} else {
		statusUrl = fmt.Sprintf("%s/job/%s/lastBuild/wfapi/describe", Config.Jenkins.Url, job)
	}
	req, err := http.NewRequest("GET", statusUrl, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(Config.Jenkins.Username, Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok: %d, resp body: %s", resp.StatusCode, resp.Body)
	}
	var build LastBuild
	err = json.NewDecoder(resp.Body).Decode(&build)
	if err != nil {
		return nil, err
	}
	return &build, nil
}

func getJenkinsUrl(job string) string {
	return fmt.Sprintf("%s/job/%s\n", Config.Jenkins.Url, job)
}

func boolToSymbol(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}
