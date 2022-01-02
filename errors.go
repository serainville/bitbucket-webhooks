package bitbucket

import "errors"

var (
	ErrInvalidSignature   = errors.New("invalid secret or digest")
	ErrMissingSecret      = errors.New("expected secret to be set")
	ErrEventType          = errors.New("invalid event type")
	ErrReadingRequestBody = errors.New("unable to read request body")
)
