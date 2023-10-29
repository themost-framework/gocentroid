package query

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FormatNameReference(name string) string {
	re := regexp.MustCompile(`^(\w+)`)
	return string(re.ReplaceAll([]byte(name), []byte("$$$1")))
}

func TrimNameReference(name string) string {
	re := regexp.MustCompile(`^\$(\w+)`)
	return string(re.ReplaceAll([]byte(name), []byte("$1")))
}

type SqlUtils struct{}

func (utils *SqlUtils) Escape(value any) string {
	if value == nil {
		return "NULL"
	}
	switch value.(type) {
	case bool:
		if value == true {
			return "true"
		} else {
			return "false"
		}
	case int:
	case int64:
	case int32:
		val := value.(int64)
		return strconv.FormatInt(val, 10)
	case float32:
		val := value.(float64)
		return strconv.FormatFloat(val, 'f', -1, 32)
	case float64:
		val := value.(float64)
		return strconv.FormatFloat(val, 'f', -1, 64)
	case time.Time:
		val := value.(time.Time)
		return val.Format("2006-01-02 15:04:05.000-07:00")
	}

	str := fmt.Sprintf("%v", value)
	return utils.EscapeString(str)

}

func (utils *SqlUtils) EscapeString(val string) string {
	res := val
	res = strings.ReplaceAll(res, "'", "\\'")
	res = strings.ReplaceAll(res, "\"", "\\\"")
	res = strings.ReplaceAll(res, "\x00", "\\0")
	res = strings.ReplaceAll(res, "\n", "\\n")
	res = strings.ReplaceAll(res, "\r", "\\r")
	res = strings.ReplaceAll(res, "\b", "\\b")
	res = strings.ReplaceAll(res, "\t", "\\t")
	res = strings.ReplaceAll(res, "\x1a", "\\Z")
	res = strings.ReplaceAll(res, "\\\\", "\\\\\\\\")
	return "'" + res + "'"
}
