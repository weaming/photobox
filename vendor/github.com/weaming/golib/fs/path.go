package fs

import (
	"os"
	"strings"
)

func ExpandUser(s string) string {
	if strings.HasPrefix(s, "$HOME") {
		s = strings.Replace(s, "$HOME", os.Getenv("HOME"), 1)
	}
	if strings.HasPrefix(s, "~/") {
		s = strings.Replace(s, "~", os.Getenv("HOME"), 1)
	}
	return s
}
