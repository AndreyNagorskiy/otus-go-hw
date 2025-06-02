package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	result := make(DomainStat)
	var mutex sync.Mutex

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			emailParts := strings.SplitN(user.Email, "@", 2)
			if len(emailParts) == 2 {
				domain := strings.ToLower(emailParts[1])
				mutex.Lock()
				result[domain]++
				mutex.Unlock()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
