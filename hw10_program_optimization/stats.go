package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var parser fastjson.Parser
	scanner := bufio.NewScanner(r)
	result := DomainStat{}

	for scanner.Scan() {
		parsed, err := parser.ParseBytes(scanner.Bytes())
		if err != nil {
			return result, err
		}

		email := string(parsed.GetStringBytes("Email"))
		emailParts := strings.SplitN(email, "@", 2)
		if len(emailParts) > 1 {
			emailDomain := strings.ToLower(emailParts[1])
			domainParts := strings.SplitN(emailDomain, ".", 2)
			if len(domainParts) > 1 && domainParts[1] == domain {
				result[emailDomain]++
			}
		}
	}

	return result, nil
}
