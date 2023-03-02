package config

import "github.com/Antonov-Alexander/otus-go/final_project/types"

type StaticConfig struct {
	BaseConfig
}

func (c *StaticConfig) Init([]int) error {
	return nil
}

func (c *StaticConfig) Set(checks map[int]types.CheckConfig) {
	c.checks = checks
}
