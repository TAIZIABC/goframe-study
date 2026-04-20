// config/config.go
// YAML 配置加载
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tasks []TaskConfig `yaml:"tasks"`
}

type TaskConfig struct {
	Name    string        `yaml:"name"`
	Cron    string        `yaml:"cron"`
	Command string        `yaml:"command"`
	Timeout time.Duration `yaml:"timeout"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置: %w", err)
	}

	for i, t := range cfg.Tasks {
		if t.Name == "" {
			return nil, fmt.Errorf("任务 #%d 缺少 name", i+1)
		}
		if t.Cron == "" {
			return nil, fmt.Errorf("任务 [%s] 缺少 cron", t.Name)
		}
		if t.Command == "" {
			return nil, fmt.Errorf("任务 [%s] 缺少 command", t.Name)
		}
		if t.Timeout == 0 {
			cfg.Tasks[i].Timeout = 60 * time.Second
		}
	}

	return &cfg, nil
}
