package bitbucketwebhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"strings"
)

// Event holds the Bitbucket Webhook event type
type Event string

// Option holds a webhook option
type Option func(*Webhook)

// Webhook is used to handle Bitbucket webhook events
type Webhook struct {
	secret                string
	preserveRequestBody   bool
	disableHMACValidation bool
}

// New creates a new Webhook with default settings. The default Webhook does not set a Webhook Secret and
// will not attempt to preserve the body of a `*http.Request` after it has been read. However, options can
// be used to change the default behaviour of a new webhook.
//
// Options:
//  WithSecret("WEBHOOK_SECRET")
//  PreserveBody()
//  WithoutHMAC()
//
// WithSecret sets the webhook secret that is used as a key when validating a Bitbucket HMAC signature.
//
// PreserveBody preserves the *http.Request body after being read by a webhook.
//
// WithoutHMAC disables HMAC validation. When set to true, an incoming Bitbucket Webhook request will not be checked to ensure
// it matches its X-Hub-Signature hash. This should not be used in production environments.
//
// Example 1: creates a default webhook
//  webhook.New()
//
// Example 2: creates a webhook with a secret and presevse the *http.Request body.
//  hook := webhook.New(WithSecret("WEBHOOK_SECRET"), PreserveBody())
//
func New(options ...Option) *Webhook {
	const (
		defaultPreserveRequestBody   = false
		defaultDisableHMACValidation = false
	)

	w := &Webhook{
		preserveRequestBody:   defaultPreserveRequestBody,
		disableHMACValidation: defaultDisableHMACValidation,
	}

	for _, opt := range options {
		opt(w)
	}

	return w
}

// WithSecret is used to set a Webook's secret. This must be used when accepting requests from a Bitbucket
// webhook that has a secret set, and the value of secret must match that of the Bitbucket webhook this library is accepting
// requests from.
func WithSecret(secret string) Option {
	return func(w *Webhook) {
		w.secret = secret
	}
}

// PreserveBody is used when the body of a `*http.Request` needs to be read by other processes. By default,
// anytime the body of a `*http.Request` is read its contents are cleared.
func PreserveBody() Option {
	return func(w *Webhook) {
		w.preserveRequestBody = true
	}
}

// WithoutHMAC diables HMAC Signature validation. All incoming events should be validated using their included
// HMAC signature, when included in a X-Hub-Signature header, to verify its authenticity and integrity. By disabling
// this check an event may come from an untrusted source or have been modified on route.
func WithoutHMAC() Option {
	return func(w *Webhook) {
		w.disableHMACValidation = true
	}
}

// Parse an a Bitbucket Webhook request. The HMAC signature of the request will be validated
// when the 'X-Hub-Signature' header key is set.
func (hook *Webhook) Parse(req *http.Request) (interface{}, error) {

	event := Event(req.Header.Get("X-Event-Key"))
	if event == "" {
		return nil, fmt.Errorf("'%s' is not a valid event type", event)
	}

	fmt.Println(req.Header)
	if event == "diagnostics:ping" {
		return DiagnosticPingEvent{Test: true}, nil
	}

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil || len(payload) == 0 {
		return nil, fmt.Errorf("could not read request body: %w", err)
	}

	if err := hook.VerifySignature(payload, req.Header.Get("X-Hub-Signature"), hook.secret); err != nil {
		return nil, fmt.Errorf("could not validate signature: %w", err)
	}

	switch event {
	case "pr:opened":
		var pl PullRequestOpenedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:declined":
		var pl PullRequestDeclinedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:deleted":
		var pl PullRequestDeletedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:comment:added":
		var pl PullRequestCommentAddedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:comment:deleted":
		var pl PullRequestCommentDeletedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:comment:edited":
		var pl PullRequestCommentEditedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:reviewer:updated":
		var pl PullRequestReviewerUpdatedPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "pr:reviewer:approved":
		fallthrough
	case "pr:reviewer:unapproved":
		fallthrough
	case "pr:reviewer:needs_work":
		var pl PullRequestReviewerPayload
		err := json.Unmarshal(payload, &pl)
		return pl, err
	case "repo:refs_changed":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "repo:modified":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "repo:forked":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "repo:comment:added":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "repo:comment:edited":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "repo:comment:deleted":
		return nil, fmt.Errorf("'%s' not implemented", event)
	case "mirror:repo_synchronized":
		return nil, fmt.Errorf("'%s' not implemented", event)
	default:
		return nil, fmt.Errorf("'%s' is not a valid Bitbucket Webhook event type", event)
	}
}

// VerifySignature is used to check an HMAC signature
func (hook *Webhook) VerifySignature(payload []byte, encodedHash, secret string) error {
	if encodedHash == "" {
		return nil
	}

	if secret == "" && encodedHash != "" {
		return errors.New("requires webhook secret to be set")
	}

	if len(payload) == 0 {
		return errors.New("payload cannot be empty")
	}

	var hashFn func() hash.Hash
	var messageMAC string

	if strings.HasPrefix(encodedHash, "sha256=") {
		messageMAC = strings.TrimPrefix(encodedHash, "sha256=")
		hashFn = sha256.New
	} else {
		prefix := strings.Split(encodedHash, "=")[0]
		return fmt.Errorf("invalid hash prefix. Expected 'sha256=...', but got: %s", prefix)
	}

	messageMACBuf, err := hex.DecodeString(messageMAC)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}

	mac := hmac.New(hashFn, []byte(secret))
	_, err = mac.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to write message as a MAC: %w", err)
	}

	expectedMAC := mac.Sum(nil)

	if ok := hmac.Equal(messageMACBuf, expectedMAC); !ok {
		return errors.New("HMAC signatures do not match")
	}

	return nil
}
