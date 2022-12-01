package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	text := "1234567890"

	t.Run("read all", func(t *testing.T) {
		reader := strings.NewReader(text)
		writer := bytes.NewBufferString("")

		transfer(reader, writer, len(text), nil)
		require.Equal(t, writer.String(), text)
	})

	t.Run("read limit", func(t *testing.T) {
		reader := strings.NewReader(text)
		writer := bytes.NewBufferString("")

		transfer(reader, writer, 3, nil)
		require.Equal(t, writer.String(), "123")
	})
}
