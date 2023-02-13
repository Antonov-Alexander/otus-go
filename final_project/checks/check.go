package checks

import (
	"errors"

	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

const (
	IpCheckType       = 1
	LoginCheckType    = 2
	PasswordCheckType = 3
)

func GetCheck(checkType int) (types.Check, error) {
	var getItem func(request types.Request) types.Item
	switch checkType {
	case IpCheckType:
		getItem = func(request types.Request) types.Item {
			return request.IP
		}
	case LoginCheckType:
		getItem = func(request types.Request) types.Item {
			return request.Login
		}
	case PasswordCheckType:
		getItem = func(request types.Request) types.Item {
			return request.Password
		}
	}

	if getItem != nil {
		return &Check{getItem: getItem}, nil
	}

	return nil, errors.New("undefined check type")
}

type Check struct {
	config  types.CheckConfig
	storage types.Storage
	getItem func(request types.Request) types.Item
}

func (b *Check) Init(config types.CheckConfig, storage types.Storage) error {
	b.config = config
	b.storage = storage
	return nil
}

func (b *Check) GetItem(request types.Request) types.Item {
	return b.getItem(request)
}

func (*Check) GetDefaultConfig() types.CheckConfig {
	return types.CheckConfig{
		BlackList:  map[types.Item]struct{}{},
		WhiteList:  map[types.Item]struct{}{},
		ItemLimits: map[types.Item][]types.Limit{},
		CommonLimits: []types.Limit{
			{
				Interval: 60,
				Limit:    100,
			},
		},
	}
}

func (b *Check) Check(request types.Request) error {
	item := b.GetItem(request)

	if _, ok := b.config.WhiteList[item]; ok {
		return nil
	}

	if _, ok := b.config.BlackList[item]; ok {
		return errors.New("blacklisted")
	}

	var limits []types.Limit
	itemLimits, ok := b.config.ItemLimits[item]
	if ok {
		limits = itemLimits
	} else {
		limits = b.config.CommonLimits
	}

	for _, limit := range limits {
		if !b.storage.Inc(item, limit) {
			return errors.New("limited")
		}
	}

	return nil
}

func (b *Check) ClearCounter(item types.Item) {
	b.storage.Reset(item)
}

func (b *Check) AddWhiteListItem(item types.Item) {
	b.config.WhiteList[item] = struct{}{}
}

func (b *Check) AddBlackListItem(item types.Item) {
	b.config.BlackList[item] = struct{}{}
}

func (b *Check) RemoveWhiteListItem(item types.Item) {
	if _, ok := b.config.WhiteList[item]; ok {
		delete(b.config.WhiteList, item)
	}
}

func (b *Check) RemoveBlackListItem(item types.Item) {
	if _, ok := b.config.BlackList[item]; ok {
		delete(b.config.BlackList, item)
	}
}
