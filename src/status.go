package jk

import (
	"fmt"

	"github.com/schollz/progressbar/v3"
)

func Status() error {
	p, err := GetParam(false)
	if err != nil {
		return err
	}

	stageCount, err := getStageCount(p.Job)
	if err != nil {
		return err
	}
	bar := progressbar.NewOptions(
		stageCount,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionUseANSICodes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	inProgress := make(map[string]bool, 0)
	succeed := make(map[string]bool, 0)
	for {
		build, err := getLastBuild(p.Job, false)
		if err != nil {
			return err
		}
		for _, stage := range build.Stages {
			name := stage.Name
			switch stage.Status {
			case stageStatusSuccess:
				if succeed[name] {
					continue
				}
				succeed[name] = true
				err = bar.Add(1)
				if err != nil {
					return err
				}
			case stageStatusInProgress:
				if inProgress[name] {
					continue
				}
				inProgress[name] = true
				bar.Describe(name)
			case stageStatusAborted:
				return fmt.Errorf("\ndeploy aborted, please check: %s\n", getJenkinsUrl(p.Job))
			case stageStatusFailed:
				return fmt.Errorf("\ndeploy failed, please check: %s\n", getJenkinsUrl(p.Job))
			default:
				panic(fmt.Errorf("unexpected stage status: %s", stage.Status))
			}
		}
		if len(succeed) == stageCount {
			break
		}
		sleep()
	}
	return nil
}

func getStageCount(job string) (int, error) {
	build, err := getLastBuild(job, true)
	if err != nil {
		return -1, err
	}
	return len(build.Stages), nil
}
