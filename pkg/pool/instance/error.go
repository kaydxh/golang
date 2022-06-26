package instance

type ContextCanceledError struct {
	Message string
}

func (e ContextCanceledError) Error() string {
	return e.Message
}

type TimeoutError struct {
	Message string
}

func (e TimeoutError) Error() string {
	return e.Message
}
