package errs

type ErrBadRequest struct {
	message string
}

func NewBadRequest(message string) ErrBadRequest {
	return ErrBadRequest{message: message}
}

func (err ErrBadRequest) Error() string {
	if err.message == "" {
		return "bad request"

	} else {
		return err.message
	}
}
