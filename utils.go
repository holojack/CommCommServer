package main

import (
	"errors"
	"strings"
)

func splitString(s string, delim string) ([]string, error) {
	if strings.Contains(s, delim) {
		return strings.Split(s, delim), nil
	}
	return nil, errors.New("Substring not found")
}
