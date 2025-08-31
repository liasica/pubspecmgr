// Copyright (C) pubspecmgr. 2025-present.
//
// Created at 2025-08-30, by liasica

package pubspecmgr

import (
	_ "embed"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
)

//go:embed pubspecmgr.yaml
var defaultConfig []byte

var (
	config *Config
	once   sync.Once
)

type Config struct {
	MarkedPaths map[string]map[string]string `yaml:"marked_paths"`

	marked []*yaml.Path
}

func init() {
	logger := log.New(os.Stderr)
	log.SetDefault(logger)
}

// LoadConfig loads config from file or use default config
func LoadConfig(params ...string) {
	once.Do(func() {
		b := defaultConfig
		var err error
		if len(params) > 0 && params[0] != "" {
			b, err = os.ReadFile(params[0])
			if err != nil {
				log.Fatalf("failed to read config file: %v", err)
			}
		}
		config = new(Config)
		err = yaml.Unmarshal(b, config)
		if err != nil {
			log.Fatalf("failed to unmarshal config: %v", err)
		}
	})
}

func GetConfig() *Config {
	if config == nil {
		LoadConfig()
	}
	return config
}

func (c *Config) GetMarked() []*yaml.Path {
	if c.marked != nil {
		return c.marked
	}

	var marked []*yaml.Path
	for p, m := range c.MarkedPaths {
		for k := range m {
			builder := &yaml.PathBuilder{}
			builder = builder.Root().Child(p)
			builder.Child(k)

			marked = append(marked, builder.Build())
		}
	}

	c.marked = marked
	return marked
}
