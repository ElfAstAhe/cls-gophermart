package loyalty

import "fmt"

type LSClientError struct {
	message string
	err     error
}

var LSClientErr *LSClientError

func NewLSClientError(message string, err error) *LSClientError {
	return &LSClientError{
		message: message,
		err:     err,
	}
}

func (lc *LSClientError) Error() string {
	return fmt.Sprintf("loyalty system client error [%v] with message [%s]", lc.err, lc.message)
}

func (lc *LSClientError) Unwrap() error {
	return lc.err
}
