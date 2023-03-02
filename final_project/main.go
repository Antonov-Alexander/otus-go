package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Antonov-Alexander/otus-go/final_project/checker"
	"github.com/Antonov-Alexander/otus-go/final_project/checks"
	"github.com/Antonov-Alexander/otus-go/final_project/config"
	"github.com/Antonov-Alexander/otus-go/final_project/storages"
	"github.com/Antonov-Alexander/otus-go/final_project/types"
)

const (
	CheckRoute             = "/check"
	ResetRoute             = "/reset"
	AddIpWhiteListRoute    = "/whitelist/ip/add"
	AddIpBlackListRoute    = "/blacklist/ip/add"
	RemoveIpWhiteListRoute = "/whitelist/ip/remove"
	RemoveIpBlackListRoute = "/blacklist/ip/remove"
)

func main() {
	// examples
	// curl http://localhost:5000/check?ip=12345
	// curl http://localhost:5000/reset?ip=12345
	// curl http://localhost:5000/whitelist/ip/add?ip=12345
	// curl http://localhost:5000/whitelist/ip/remove?ip=12345

	serverHost := flag.String("server_host", "localhost", "server host")
	serverPort := flag.String("server_port", "5000", "server port")
	dbHost := flag.String("db_host", "localhost", "database host")
	dbPort := flag.Int("db_port", 5432, "database port")
	dbName := flag.String("db_name", "postgres", "database name")
	dbUser := flag.String("db_user", "admin", "database user")
	dbPass := flag.String("db_pass", "admin", "database password")
	flag.Parse()

	checkerChecker, err := initChecker(*dbHost, *dbPort, *dbName, *dbUser, *dbPass)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = initServer(checkerChecker, *serverHost, *serverPort); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initChecker(host string, port int, db, user, password string) (checker.Checker, error) {
	checkerConfig := &config.DbConfig{
		Host:     host,
		Port:     port,
		DbName:   db,
		User:     user,
		Password: password,
	}

	checkerStorageType := storages.MemoryStorageType
	checkTypes := []int{
		checks.IpCheckType,
		checks.LoginCheckType,
		checks.PasswordCheckType,
	}

	checkerChecker := checker.Checker{}
	if err := checkerChecker.Init(checkTypes, checkerStorageType, checkerConfig); err != nil {
		return checkerChecker, err
	}

	return checkerChecker, nil
}

func initServer(checkerChecker checker.Checker, host, port string) error {
	http.HandleFunc(CheckRoute, getCheckFunc(checkerChecker))
	http.HandleFunc(ResetRoute, getResetFunc(checkerChecker))

	listRoutes := []string{
		AddIpWhiteListRoute,
		AddIpBlackListRoute,
		RemoveIpWhiteListRoute,
		RemoveIpBlackListRoute,
	}

	for _, listRoute := range listRoutes {
		http.HandleFunc(listRoute, getListFunc(listRoute, checkerChecker))
	}

	addr := host + ":" + port
	fmt.Println("Starting server on", addr)

	err := http.ListenAndServe(addr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		return errors.New("server closed")
	} else if err != nil {
		return errors.New(fmt.Sprintf("error starting server: %s\n", err))
	}

	return nil
}

func parseRequestFromQuery(query url.Values) types.Request {
	var ip int
	if value, err := strconv.Atoi(query.Get("ip")); err == nil {
		ip = value
	}

	return types.Request{
		IP:       ip,
		Login:    query.Get("login"),
		Password: query.Get("password"),
	}
}

func getCheckFunc(checkerChecker checker.Checker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := parseRequestFromQuery(r.URL.Query())
		result := "0"
		defer func() {
			_, _ = io.WriteString(w, result)
		}()

		if err := checkerChecker.Check(request); err != nil {
			fmt.Println("check_err:", err)
			result = "1"
		}
	}
}

func getResetFunc(checkerChecker checker.Checker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := parseRequestFromQuery(r.URL.Query())
		result := "1"
		defer func() {
			_, _ = io.WriteString(w, result)
		}()

		checkTypes := []int{
			checks.IpCheckType,
			checks.LoginCheckType,
		}

		for checkType := range checkTypes {
			check, err := checks.GetCheck(checkType)
			if err != nil {
				continue
			}

			item := check.GetItem(request)
			if item == nil {
				continue
			}

			if err = checkerChecker.ClearCounter(checkType, request); err != nil {
				fmt.Println("reset_err:", err)
				result = "0"
			}
		}
	}
}

func getListFunc(route string, checkerChecker checker.Checker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		request := parseRequestFromQuery(r.URL.Query())
		result := "1"
		defer func() {
			_, _ = io.WriteString(w, result)
		}()

		var checkType int
		var checkerMethod string

		switch route {
		case AddIpWhiteListRoute:
			checkType = checks.IpCheckType
			checkerMethod = checker.AddWhiteListItemMethod
		case AddIpBlackListRoute:
			checkType = checks.IpCheckType
			checkerMethod = checker.AddBlackListItemMethod
		case RemoveIpWhiteListRoute:
			checkType = checks.IpCheckType
			checkerMethod = checker.RemoveWhiteListItemMethod
		case RemoveIpBlackListRoute:
			checkType = checks.IpCheckType
			checkerMethod = checker.RemoveBlackListItemMethod
		}

		if err := checkerChecker.CallListMethod(checkType, request, checkerMethod); err != nil {
			fmt.Println("check_err:", err)
			result = "0"
		}
	}
}
