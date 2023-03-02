package config

import (
	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

type BaseConfig struct {
	checks map[int]types.CheckConfig
}

func (c *BaseConfig) Init([]int) error {
	c.checks = map[int]types.CheckConfig{}
	return nil
}

func (c *BaseConfig) GetCheckConfig(checkName int) (types.CheckConfig, bool) {
	result, ok := c.checks[checkName]
	return result, ok
}

func (*BaseConfig) getNewCheckConfig() types.CheckConfig {
	return types.CheckConfig{
		ItemLimits: map[types.Item][]types.Limit{},
		BlackList:  map[types.Item]struct{}{},
		WhiteList:  map[types.Item]struct{}{},
	}
}

func (*BaseConfig) AddWhiteListItem(int, types.Request) error    { return nil }
func (*BaseConfig) AddBlackListItem(int, types.Request) error    { return nil }
func (*BaseConfig) RemoveWhiteListItem(int, types.Request) error { return nil }
func (*BaseConfig) RemoveBlackListItem(int, types.Request) error { return nil }
