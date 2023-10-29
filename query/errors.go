package query

type NotImplementedError struct{}

func (m *NotImplementedError) Error() string {
	return "not yet implemented"
}

type UnknowMethodError struct{}

func (m *UnknowMethodError) Error() string {
	return "the specified method is unknown"
}
