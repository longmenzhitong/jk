package jk

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultPollingIntervalSecond = 1
)

const configFileName = "config.yml"

var Config config

type config struct {
	Jenkins struct {
		Url                   string `yaml:"url"`
		Username              string `yaml:"username"`
		Password              string `yaml:"password"`
		PrintStatus           bool   `yaml:"printStatus"`
		PollingIntervalSecond int    `yaml:"pollingIntervalSecond"`
	} `yaml:"jenkins"`
}

func (c *config) Init() {
	dir := ProjectDir()
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	configPath := Path(configFileName)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		Config = config{}
		Config.check()
		out, err := yaml.Marshal(Config)
		if err != nil {
			panic(err)
		}
		RewriteLinesToPath(configPath, []string{string(out)})
		return
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	c.check()
}

func (c *config) check() {
	if c.Jenkins.PollingIntervalSecond <= 0 {
		c.Jenkins.PollingIntervalSecond = defaultPollingIntervalSecond
	}
}
