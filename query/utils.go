package query

import "regexp"

func FormatNameReference(name string) string {
	re := regexp.MustCompile(`^(\w+)`)
	return string(re.ReplaceAll([]byte(name), []byte("$$$1")))
}

func TrimNameReference(name string) string {
	re := regexp.MustCompile(`^\$(\w+)`)
	return string(re.ReplaceAll([]byte(name), []byte("$1")))
}
