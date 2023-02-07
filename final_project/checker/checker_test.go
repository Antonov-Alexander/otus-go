package checker

import (
	"testing"

	"github.com/Antonov-Alexander/otus-go/final_project/checks"
	"github.com/Antonov-Alexander/otus-go/final_project/config"
	"github.com/Antonov-Alexander/otus-go/final_project/storages"
	"github.com/Antonov-Alexander/otus-go/final_project/types"
	"github.com/stretchr/testify/require"
)

func TestChecker(t *testing.T) {
	t.Run("single item checks", func(t *testing.T) {
		testData := []struct {
			CheckType int
			Request   types.Request
		}{
			{
				CheckType: checks.IpCheckType,
				Request:   types.Request{IP: 12345},
			},
			{
				CheckType: checks.LoginCheckType,
				Request:   types.Request{Login: "TestLogin"},
			},
		}

		storageType := storages.MemoryStorageType

		for _, testItem := range testData {
			testTypes := []int{testItem.CheckType}
			testConfig := &config.StaticConfig{}
			testConfig.Set(map[int]types.CheckConfig{
				testItem.CheckType: {
					CommonLimits: []types.Limit{
						{
							Interval: 2,
							Limit:    4,
						},
						{
							Interval: 1,
							Limit:    2,
						},
					},
				},
			})

			checkerChecker := Checker{}
			err := checkerChecker.Init(testTypes, storageType, testConfig)
			require.NoError(t, err)

			require.NoError(t, checkerChecker.Check(testItem.Request))
			require.NoError(t, checkerChecker.Check(testItem.Request))
			require.Error(t, checkerChecker.Check(testItem.Request))
		}
	})

	t.Run("white list checks", func(t *testing.T) {
		testData := []struct {
			CheckType int
			Request   types.Request
		}{
			{
				CheckType: checks.IpCheckType,
				Request:   types.Request{IP: 12345},
			},
			{
				CheckType: checks.LoginCheckType,
				Request:   types.Request{Login: "TestLogin"},
			},
		}

		storageType := storages.MemoryStorageType

		for _, testItem := range testData {
			testTypes := []int{testItem.CheckType}
			testConfig := &config.StaticConfig{}
			testConfig.Set(map[int]types.CheckConfig{
				testItem.CheckType: {
					ItemLimits: map[types.Item][]types.Limit{},
					BlackList:  map[types.Item]struct{}{},
					WhiteList:  map[types.Item]struct{}{},
					CommonLimits: []types.Limit{
						{
							Interval: 1,
							Limit:    1,
						},
					},
				},
			})

			checkerChecker := Checker{}
			err := checkerChecker.Init(testTypes, storageType, testConfig)
			require.NoError(t, err)

			// добавим в whitelist
			err = checkerChecker.AddWhiteListItem(testItem.CheckType, testItem.Request)
			require.NoError(t, err)

			require.NoError(t, checkerChecker.Check(testItem.Request))
			require.NoError(t, checkerChecker.Check(testItem.Request))

			// уберём из whitelist-а
			err = checkerChecker.RemoveWhiteListItem(testItem.CheckType, testItem.Request)
			require.NoError(t, err)

			require.NoError(t, checkerChecker.Check(testItem.Request))
			require.Error(t, checkerChecker.Check(testItem.Request))
		}
	})

	t.Run("black list checks", func(t *testing.T) {
		testData := []struct {
			CheckType int
			Request   types.Request
		}{
			{
				CheckType: checks.IpCheckType,
				Request:   types.Request{IP: 12345},
			},
			{
				CheckType: checks.LoginCheckType,
				Request:   types.Request{Login: "TestLogin"},
			},
		}

		storageType := storages.MemoryStorageType

		for _, testItem := range testData {
			testTypes := []int{testItem.CheckType}
			testConfig := &config.StaticConfig{}
			testConfig.Set(map[int]types.CheckConfig{
				testItem.CheckType: {
					ItemLimits: map[types.Item][]types.Limit{},
					BlackList:  map[types.Item]struct{}{},
					WhiteList:  map[types.Item]struct{}{},
					CommonLimits: []types.Limit{
						{
							Interval: 1,
							Limit:    1,
						},
					},
				},
			})

			checkerChecker := Checker{}
			err := checkerChecker.Init(testTypes, storageType, testConfig)
			require.NoError(t, err)

			// добавим в blacklist
			err = checkerChecker.AddBlackListItem(testItem.CheckType, testItem.Request)
			require.NoError(t, err)

			require.Error(t, checkerChecker.Check(testItem.Request))
			require.Error(t, checkerChecker.Check(testItem.Request))

			// уберём из blacklist-а
			err = checkerChecker.RemoveBlackListItem(testItem.CheckType, testItem.Request)
			require.NoError(t, err)

			require.NoError(t, checkerChecker.Check(testItem.Request))
			require.Error(t, checkerChecker.Check(testItem.Request))
		}
	})
}
