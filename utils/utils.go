package utils

import (
	"regexp"
	"strings"
)

func FormatKey(text string) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	text = strings.ToUpper(text)
	reg, _ := regexp.Compile(`[^A-Z0-9]+`)
	processedString := reg.ReplaceAllString(text, "_")
	return processedString
}

func CheckSliceHasTheElement(commaSeperatedString []string, element string) bool {
	found := false
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for _, ele := range commaSeperatedString {
		if strings.EqualFold(ele, element) {
			found = true
			break
		}
	}
	return found
}

func FormatMessage(text string) string {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	processedString := strings.ReplaceAll(text, ":", "-")
	return processedString
}
