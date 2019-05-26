package utils

import (
	"fmt"
)

var okayResponses = []string{"y", "Y", "yes", "Yes", "YES"}

func Confirm(msg string) bool {
	fmt.Printf("%s ", msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	for _, r := range okayResponses {
		if r == response {
			return true
		}
	}

	return false
}
