package query

import "errors"

type SqlFormatter struct {
	NameFormat string
}

func (formatter *SqlFormatter) escapeName(element QueryElement) *string {
	if element["$name"] != nil {
	}
	errors.New("The expression has not yet implemented")
	return nil
}

func (formatter *SqlFormatter) escape(value interface{}) *string {
	return nil
}

func (formatter *SqlFormatter) formatSelect(query QueryExpression) *string {
	return nil
}
