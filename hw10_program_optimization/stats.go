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
		for pos, char := range email {
			if char == 64 {
				emailDomain := strings.ToLower(email[pos+1:])
				for domainPos, domainChar := range emailDomain {
					if domainChar == 46 && domain == strings.ToLower(emailDomain[domainPos+1:]) {
						result[emailDomain]++
					}
				}
			}
		}
	}

	return result, nil
}
