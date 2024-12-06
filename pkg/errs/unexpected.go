package errs

import "log/slog"

type ErrUnexpected struct {
	message string
}

func NewErrUnexpected(err error) ErrUnexpected {
	return ErrUnexpected{
		message: err.Error(),
	}
}

func (err ErrUnexpected) Error() string {
	slog.Error(err.message)

	return "unexpected error"
}
