package main

const (
	OK = 0

	ERROR_NOTMATCH = 5001

	ERROR_UNMARSHAL = 8001
	ERROR_MARSHAL   = 8002
)

type Error struct {
	code    int
	message string
}

func (e *Error) String() string {
	return e.message
}

func (e *Error) Code() int {
	return e.code
}
