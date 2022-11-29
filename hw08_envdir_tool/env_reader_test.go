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
	}

	testParams := struct {
		dir        string
		customVars CustomVars
	}{
		dir: "testdata/test/",
		customVars: CustomVars{
			{name: "TEST_1", value: "TEST_VALUE_1", needRemove: false},
			{name: "TEST_2", value: "", needRemove: true},
		},
	}

	// init
	_ = os.Mkdir(testParams.dir, 0o755)
	for _, customVar := range testParams.customVars {
		_ = os.WriteFile(testParams.dir+customVar.name, []byte(customVar.value), 0o755)
	}

	// cleanup
	t.Cleanup(func() {
		for _, customVar := range testParams.customVars {
			_ = os.Remove(testParams.dir + customVar.name)
		}
		_ = os.Remove(testParams.dir)
	})

	// test
	t.Run("Common test", func(t *testing.T) {
		if env, err := ReadDir(testParams.dir); err == nil {
			for _, customVar := range testParams.customVars {
				value, ok := env[customVar.name]
				require.True(t, ok)
				require.Equal(t, value.Value, customVar.value)
				require.Equal(t, value.NeedRemove, customVar.needRemove)
			}

			_, ok := env["UNDEFINED"]
			require.False(t, ok)
			require.Equal(t, len(testParams.customVars), len(env))
		}
	})
}
