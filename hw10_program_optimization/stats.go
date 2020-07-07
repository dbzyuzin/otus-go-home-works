package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/buger/jsonparser"
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

var (
	ErrDataCorrupted = errors.New("data is corrupted")
	ErrEmptyDomain   = errors.New("empty domain")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}

	stat := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		email, err := jsonparser.GetString(scanner.Bytes(), "Email")
		if err != nil {
			return nil, ErrDataCorrupted
		}
		if email != "" && strings.HasSuffix(email, domain) {
			eDomain := strings.SplitN(email, "@", 2)[1]
			stat[strings.ToLower(eDomain)]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stat, nil
}
