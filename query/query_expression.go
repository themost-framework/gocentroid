package query

type QueryElement map[string]interface{}

type QueryExpression struct {
	Query QueryElement
}

func (query *QueryExpression) from(collection string) *QueryExpression {
	if query.Query == nil {
		query.Query = make(QueryElement)
	}
	query.Query["$collection"] = QueryElement{
		collection: 1,
	}
	return query
}

func (query *QueryExpression) pick(attr ...string) *QueryExpression {
	if query.Query == nil {
		query.Query = make(QueryElement)
	}
	attributes := QueryElement{}
	for _, attribute := range attr {
		attributes[attribute] = 1
	}
	query.Query["$select"] = attributes
	return query
}
