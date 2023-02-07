package main

import (
	"log"

	"github.com/Antonov-Alexander/otus-go/final_project/checker"
	"github.com/Antonov-Alexander/otus-go/final_project/checks"
	"github.com/Antonov-Alexander/otus-go/final_project/config"
	"github.com/Antonov-Alexander/otus-go/final_project/storages"
)

func main() {
	// TODO сделать Config, который подгружает конфиг из базы
	checkerConfig := &config.BaseConfig{}

	checkerStorageType := storages.MemoryStorageType
	checkTypes := []int{
		checks.IpCheckType,
		checks.LoginCheckType,
		checks.PasswordCheckType,
	}

	checkerChecker := checker.Checker{}
	if err := checkerChecker.Init(checkTypes, checkerStorageType, checkerConfig); err != nil {
		log.Fatalf("initializing error: %s", err)
	}

	// апдейт параметров

	// сервер

}
