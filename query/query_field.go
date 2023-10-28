package query

type QueryField map[string]interface{}

func (field *QueryField) from(collection string) *QueryField {
	return field
}
