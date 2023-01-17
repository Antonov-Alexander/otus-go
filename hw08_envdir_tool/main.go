package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) < 2 {
		return
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
