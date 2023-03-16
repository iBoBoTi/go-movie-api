package swapi

type Error struct {
	Message     string
	StatusCode  int
	ActualError error
}

func (e Error) Error() string {
	return e.Message
}
