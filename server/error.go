package nc

// New returns an error that formats as the given text.

func New(text string) error {
	return &Error{text}
}

// Error is a trivial implementation of error.
type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

func (e *Error) String() string {
	return e.s
}
