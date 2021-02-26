package config

import (
	"fmt"
	"regexp"
	"strings"
)

func CamelCase(str string) (string, error) {
	if regexp.MustCompile("[^a-zA-Z0-9_]+").MatchString(str) {
		return "", fmt.Errorf("config variables can only contain chars: [a-zA-Z0-9_]")
	}
	link := regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
	str = strings.ToLower(str)
	converted := link.ReplaceAllStringFunc(str,
		func(s string) string {
			return strings.ToUpper(strings.Replace(s, "_", "", -1))
		})
	return converted, nil
}
