package jk

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	stageStatusSuccess    = "SUCCESS"
	stageStatusInProgress = "IN_PROGRESS"
	stageStatusAborted    = "ABORTED"
	stageStatusFailed     = "FAILED"
)

func checkConfig() error {
	PrintStartMsg("Check config...\n")
	rows := make([]table.Row, 0)

	if Config.Jenkins.Url == "" {
		return &InvalidConfigError{Message: "jenkins.url doesn't exist"}
	}
	rows = append(rows, table.Row{"Url", boolToSymbol(true)})

	if Config.Jenkins.Username == "" {
		return &InvalidConfigError{Message: "jenkins.username doesn't exist"}
	}
	rows = append(rows, table.Row{"Username", boolToSymbol(true)})

	if Config.Jenkins.Password == "" {
		return &InvalidConfigError{Message: "jenkins.password doesn't exist"}
	}
	rows = append(rows, table.Row{"Password", boolToSymbol(true)})

	PrintTable(nil, rows)

	PrintDoneMsg("Check done.\n")
	return nil
}

func checkParam() (*param, error) {
	PrintStartMsg("Check params...\n")
	rows := make([]table.Row, 0)

	p, err := GetParam(true)
	if err != nil {
		return nil, err
	}

	rows = append(rows, table.Row{"JOB", p.Job})
	for k, v := range p.BuildParams {
		rows = append(rows, table.Row{k, v})
	}

	PrintTable(nil, rows)

	PrintDoneMsg("Check done.\n")
	return p, nil
}

func Deploy() error {
	// 检查配置
	err := checkConfig()
	if err != nil {
		return err
	}
	fmt.Println()

	// 检查参数
	p, err := checkParam()
	if err != nil {
		return err
	}
	fmt.Println()

	// 确认构建
	fmt.Println("Sure to deploy? [y/n]")
	input := strings.ToLower(ReadString())
	if input != "y" {
		return &CancelledError{}
	}

	// 发起构建
	err = startBuild(p)
	if err != nil {
		return err
	}

	if Config.Jenkins.PrintStatus {
		// 等待构建
		err = waitDeploy(p.Job)
		if err != nil {
			return err
		}

		// 打印状态
		err = Status()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("To check the progress, please visit: %s", getJenkinsUrl(p.Job))
	}

	fmt.Println("🍺 Deploy done!")
	return nil
}

func startBuild(p *param) error {
	buildUrl := fmt.Sprintf("%s/job/%s/buildWithParameters", Config.Jenkins.Url, p.Job)
	requestBody := url.Values{}
	for k, v := range p.BuildParams {
		requestBody.Add(k, v)
	}
	req, err := http.NewRequest("POST", buildUrl, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.SetBasicAuth(Config.Jenkins.Username, Config.Jenkins.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("http status not created: %d, resp body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func waitDeploy(job string) error {
	fmt.Println("Wait deploy...")
	for {
		build, err := getLastBuild(job, false)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			if stage.Status == stageStatusInProgress {
				fmt.Println("Start deploy...")
				return nil
			}
		}
		time.Sleep(time.Duration(Config.Jenkins.PollingIntervalSecond) * time.Second)
	}
}
