package query

type NotImplementedError struct{}

func (m *NotImplementedError) Error() string {
	return "not yet implemented"
}

type UnknownMethodError struct{}

func (m *UnknownMethodError) Error() string {
	return "the specified method is unknown"
}

type InvalidExpressionError struct{}

func (m *InvalidExpressionError) Error() string {
	return "the specified expression is invalid"
}
