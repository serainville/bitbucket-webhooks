package bitbucket

import "errors"

var (
	// ErrInvalidSignature is used when an HMAC signature cannot be valdiated
	ErrInvalidSignature = errors.New("invalid secret or digest")
	// ErrMissingSecret is used when a webhook's secret is not set
	ErrMissingSecret = errors.New("expected secret to be set")
	// ErrEventType is used when an incoming eventKey does not match any known keys
	ErrEventType = errors.New("invalid event type")
	// ErrReadingRequestBody is used when the request body cannot be read
	ErrReadingRequestBody = errors.New("unable to read request body")
)
