package config

import (
	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

type BaseConfig struct {
	checks map[int]types.CheckConfig
}

func (c *BaseConfig) Init() error {
	c.checks = map[int]types.CheckConfig{}
	return nil
}

func (c *BaseConfig) GetCheckConfig(checkName int) (types.CheckConfig, bool) {
	result, ok := c.checks[checkName]
	return result, ok
}
