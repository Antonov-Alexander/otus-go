package checker

import (
	"errors"

	"github.com/Antonov-Alexander/otus-go/final_project/checks"
	"github.com/Antonov-Alexander/otus-go/final_project/storages"
	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

type Checker struct {
	checks map[int]types.Check
}

const (
	AddWhiteListItemMethod    = "AddWhiteListItem"
	AddBlackListItemMethod    = "AddBlackListItem"
	RemoveWhiteListItemMethod = "RemoveWhiteListItem"
	RemoveBlackListItemMethod = "RemoveBlackListItem"
)

func (c *Checker) Init(checkTypes []int, storageType int, config types.Config) error {
	if err := config.Init(checkTypes); err != nil {
		return errors.New("config loading error")
	}

	c.checks = make(map[int]types.Check, len(checkTypes))
	for _, checkType := range checkTypes {
		check, err := checks.GetCheck(checkType)
		if err != nil {
			return err
		}

		checkConfig, ok := config.GetCheckConfig(checkType)
		if !ok {
			checkConfig = check.GetDefaultConfig()
		}

		storage, err := storages.GetStorage(storageType)
		if err != nil {
			return err
		}

		if err = check.Init(checkConfig, storage); err != nil {
			return errors.New("check initializing error")
		}

		c.checks[checkType] = check
	}

	return nil
}

func (c *Checker) Check(request types.Request) (err error) {
	// возвращаем первую ошибку, но делаем все проверки, чтобы увеличились счётчики
	for _, check := range c.checks {
		if checkErr := check.Check(request); err == nil && checkErr != nil {
			err = checkErr
		}
	}

	return err
}

func (c *Checker) CallListMethod(checkType int, request types.Request, method string) error {
	if check, ok := c.checks[checkType]; ok {
		switch method {
		case AddWhiteListItemMethod:
			c.checks[checkType].AddWhiteListItem(check.GetItem(request))
		case AddBlackListItemMethod:
			c.checks[checkType].AddBlackListItem(check.GetItem(request))
		case RemoveWhiteListItemMethod:
			c.checks[checkType].RemoveWhiteListItem(check.GetItem(request))
		case RemoveBlackListItemMethod:
			c.checks[checkType].RemoveBlackListItem(check.GetItem(request))
		}
		return nil
	}

	return errors.New("unsupported check type")
}

func (c *Checker) AddWhiteListItem(checkType int, request types.Request) error {
	return c.CallListMethod(checkType, request, AddWhiteListItemMethod)
}

func (c *Checker) AddBlackListItem(checkType int, request types.Request) error {
	return c.CallListMethod(checkType, request, AddBlackListItemMethod)
}

func (c *Checker) RemoveWhiteListItem(checkType int, request types.Request) error {
	return c.CallListMethod(checkType, request, RemoveWhiteListItemMethod)
}

func (c *Checker) RemoveBlackListItem(checkType int, request types.Request) error {
	return c.CallListMethod(checkType, request, RemoveBlackListItemMethod)
}

func (c *Checker) ClearCounter(checkType int, request types.Request) error {
	if check, ok := c.checks[checkType]; ok {
		c.checks[checkType].ClearCounter(check.GetItem(request))
		return nil
	}

	return errors.New("unsupported check type")
}
