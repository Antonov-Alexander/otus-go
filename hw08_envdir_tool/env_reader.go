package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]EnvValue)
	zeroString := string([]byte{0x00})

	for _, dirFile := range dirFiles {
		if dirFile.IsDir() {
			continue
		}

		file, err := os.Open(dir + "/" + dirFile.Name())
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		value, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		value = strings.ReplaceAll(strings.TrimRight(value, "\n\t "), zeroString, "\n")
		result[dirFile.Name()] = EnvValue{value, value == ""}
	}

	return result, nil
}
