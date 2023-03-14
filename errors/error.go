package errors

type Error struct {
	Message string
	Status  int
}

func (e *Error) Error() string {
	return e.Message
}

func New(message string, code int) *Error {
	return &Error{
		Message: message,
		Status:  code,
	}
}
