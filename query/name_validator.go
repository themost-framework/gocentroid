package query

import (
	"errors"
	"regexp"
)

const DefaultNamePattern = "([a-zA-Z0-9_]+)"
const LatinNamePattern = "([\u0030-\u0039\u0041-\u005A\u0061-\u007A\u005F]+)"
const LatinExtendedNamePattern = "([\u0030-\u0039\u0041-\u005A\u0061-\u007A\u00A0-\u024F\u005F]+)"
const GreekNamePattern = "([\u0030-\u0039\u0041-\u005A\u0061-\u007A\u0386-\u03CE\u005F]+)"
const CyrillicNamePattern = "([\u0030-\u0039\u0041-\u007A\u0061-\u007A\u0400-\u04FF\u005F]+)"
const HebrewNamePattern = "([\u0030-\u0039\u0041-\u005A\u0061-\u007A\u05D0-\u05F2\u005F]+)"

type ObjectNameValidator struct {
	Pattern          string
	QualifiedPattern string
}

func (validator *ObjectNameValidator) Init(pattern string) *ObjectNameValidator {
	validator.Pattern = `^` + pattern + `$`
	validator.QualifiedPattern = `(^\*$)|^` + pattern + `((\.|\/)` + pattern + `)*(\.\*)?$`
	return validator
}

func (validator *ObjectNameValidator) Test(name string, qualified bool, throwError bool) (result bool, err error) {
	var matched bool
	if qualified {
		matched, _ = regexp.MatchString(validator.QualifiedPattern, name)
		if !matched && throwError {
			return matched, errors.New("invalid database object name")
		}
	} else {
		matched, _ = regexp.MatchString(validator.Pattern, name)
	}
	return matched, nil
}

func (validator *ObjectNameValidator) Escape(name string, format string) (escaped string, err error) {
	valid, err := validator.Test(name, false, true)
	if !valid {
		return "", err
	}
	re := regexp.MustCompile(validator.Pattern)
	return string(re.ReplaceAll([]byte(name), []byte(format))), nil
}

func DefaultNameValidator() *ObjectNameValidator {
	return (&ObjectNameValidator{}).Init(DefaultNamePattern)
}
