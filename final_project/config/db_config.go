package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Antonov-Alexander/otus-go/final_project/checks"
	"github.com/Antonov-Alexander/otus-go/final_project/types"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	BaseConfig
	Host     string
	Port     int
	User     string
	Password string
	DbName   string

	connection *sql.DB
}

func (c *DbConfig) Init(checkTypes []int) error {
	c.checks = map[int]types.CheckConfig{}

	if len(checkTypes) == 0 {
		return nil
	}

	if err := c.connect(); err != nil {
		return err
	}

	defer func() {
		_ = c.disconnect()
	}()

	if err := c.loadLimits(checkTypes); err != nil {
		return err
	}

	if err := c.loadLists(checkTypes); err != nil {
		return err
	}

	return nil
}

func (c *DbConfig) loadLimits(checkTypes []int) error {
	checkNames := c.getCheckNames(checkTypes)
	if len(checkNames) == 0 {
		return nil
	}

	file, err := os.ReadFile("config/db_queries/limits.sql")
	if err != nil {
		return err
	}

	query := fmt.Sprintf(string(file), strings.Join(checkNames, ","))
	rows, err := c.connection.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var name, item sql.NullString
		var interval, limit sql.NullInt32
		if err = rows.Scan(&name, &item, &interval, &limit); err != nil {
			return err
		}

		checkType, ok := checks.GetCheckType(name.String)
		if !name.Valid || !interval.Valid || !limit.Valid || !ok {
			continue
		}

		checkConfig, ok := c.checks[checkType]
		if !ok {
			checkConfig = c.getNewCheckConfig()
		}

		// достаём общие лимиты
		if item.String == "" {
			checkConfig.CommonLimits = append(checkConfig.CommonLimits, types.Limit{
				Interval: int(interval.Int32),
				Limit:    int(limit.Int32),
			})

			c.checks[checkType] = checkConfig
			continue
		}

		// достаём частные лимиты
		request, err := c.convertItemToRequest(item.String)
		if err != nil {
			return err
		}

		check, err := checks.GetCheck(checkType)
		if err != nil {
			return err
		}

		checkItem := check.GetItem(request)
		checkConfig.ItemLimits[checkItem] = append(checkConfig.ItemLimits[checkItem], types.Limit{
			Interval: int(interval.Int32),
			Limit:    int(limit.Int32),
		})

		c.checks[checkType] = checkConfig
	}

	return nil
}

func (c *DbConfig) loadLists(checkTypes []int) error {
	checkNames := c.getCheckNames(checkTypes)
	if len(checkNames) == 0 {
		return nil
	}

	file, err := os.ReadFile("config/db_queries/lists.sql")
	if err != nil {
		return err
	}

	query := fmt.Sprintf(string(file), strings.Join(checkNames, ","))
	rows, err := c.connection.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var name, item, listType sql.NullString
		if err = rows.Scan(&name, &item, &listType); err != nil {
			return err
		}

		checkType, ok := checks.GetCheckType(name.String)
		if !name.Valid || !item.Valid || !listType.Valid || !ok {
			continue
		}

		checkConfig, ok := c.checks[checkType]
		if !ok {
			checkConfig = c.getNewCheckConfig()
		}

		request, err := c.convertItemToRequest(item.String)
		if err != nil {
			return err
		}

		check, err := checks.GetCheck(checkType)
		if err != nil {
			return err
		}

		checkItem := check.GetItem(request)
		switch listType.String {
		case "white":
			checkConfig.WhiteList[checkItem] = struct{}{}
		case "black":
			checkConfig.BlackList[checkItem] = struct{}{}
		}

		c.checks[checkType] = checkConfig
	}

	return nil
}

func (c *DbConfig) connect() error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	c.connection = db
	return nil
}

func (c *DbConfig) disconnect() error {
	return c.connection.Close()
}

func (*DbConfig) getCheckNames(checkTypes []int) []string {
	checkNames := make([]string, len(checkTypes))
	for idx, checkType := range checkTypes {
		if checkName, ok := checks.GetCheckName(checkType); ok {
			checkNames[idx] = "'" + checkName + "'"
		}
	}

	return checkNames
}

func (*DbConfig) convertItemToRequest(item string) (types.Request, error) {
	var result types.Request
	err := json.Unmarshal([]byte(item), &result)
	return result, err
}
