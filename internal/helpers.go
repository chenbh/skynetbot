package bot

import (
	"errors"
	"strings"
)

func validateFile(path string) error {
	if strings.Contains(path, "/") {
		return errors.New("no special characters allowed")
	}
	return nil
}

func contains(arr []string, str string) bool {
	for _, e := range arr {
		if e == str {
			return true
		}
	}
	return false
}
