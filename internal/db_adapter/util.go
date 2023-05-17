package db_adapter

import (
	"regexp"
	"strings"
)

func clearText(text string) string {
	text = strings.ToLower(text)
	return regexp.MustCompile(`[^a-zA-Zа-яА-Я]+`).ReplaceAllString(text, "")
}
