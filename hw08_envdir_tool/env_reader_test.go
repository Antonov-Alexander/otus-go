package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	type CustomVars []struct {
		name       string
		value      string
		needRemove bool
		existed    bool
	}

	testParams := struct {
		dir        string
		customVars CustomVars
	}{
		dir: "testdata/test/",
		customVars: CustomVars{
			{
				name:       "TEST_1",
				value:      "TEST_VALUE_1",
				needRemove: false,
				existed:    true,
			},
			{
				name:       "TEST_2",
				value:      "",
				needRemove: true,
				existed:    true,
			},
			{
				name:       "TEST_=",
				value:      "TEST_VALUE_=",
				needRemove: false,
				existed:    false,
			},
		},
	}

	// init
	err := os.Mkdir(testParams.dir, 0o755)
	require.NoError(t, err)
	for _, customVar := range testParams.customVars {
		err = os.WriteFile(testParams.dir+customVar.name, []byte(customVar.value), 0o755)
		require.NoError(t, err)
	}

	// cleanup
	t.Cleanup(func() {
		for _, customVar := range testParams.customVars {
			err = os.Remove(testParams.dir + customVar.name)
			require.NoError(t, err)
		}
		err = os.Remove(testParams.dir)
		require.NoError(t, err)
	})

	// test
	t.Run("Common test", func(t *testing.T) {
		if env, err := ReadDir(testParams.dir); err == nil {
			existedCount := 0
			for _, customVar := range testParams.customVars {
				value, ok := env[customVar.name]
				if customVar.existed {
					require.True(t, ok, customVar.name)
					require.Equal(t, value.Value, customVar.value, customVar.name)
					require.Equal(t, value.NeedRemove, customVar.needRemove, customVar.name)
					existedCount++
				} else {
					require.False(t, ok, customVar.name, customVar.name)
				}
			}

			_, ok := env["UNDEFINED"]
			require.False(t, ok)
			require.Equal(t, existedCount, len(env))
		}
	})
}
