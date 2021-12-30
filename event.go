package events

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
type Option string

// Webhook is used to handle Bitbucket webhook events
type Webhook struct {
	secret string
}

// New creates a new Webhook with default settings
func New(options ...Option) *Webhook {
	return &Webhook{}
}

// Secret is used to set a Webook's secret
func (hook *Webhook) Secret(value string) {
	hook.secret = value
}

// Parse an a Bitbucket Webhook request. The HMAC signature of the request will be validated
// when the 'X-Hub-Signature' header key is set.
func (hook *Webhook) Parse(req *http.Request, events ...Event) (interface{}, error) {

	event := Event(req.Header.Get("X-Event-Key"))
	if event == "" {
		return nil, fmt.Errorf("'%s' is not a valid event type", event)
	}

	if event == "diagnostic:ping" {
		return DiagnosticPingEvent{Test: true}, nil
	}

	payload, err := ioutil.ReadAll(req.Body)
	if err != nil || len(payload) == 0 {
		return nil, fmt.Errorf("could not read request body: %w", err)
	}

	if err := hook.VerifySignature(payload, req.Header.Get("X-Hub-Signature"), hook.secret); err != nil {
		return nil, fmt.Errorf("could not validate signature: %w", err)
	}
	var bitbucketEvent Event

	var found bool
	for _, evt := range events {
		found = true
		bitbucketEvent = evt
		break
	}

	if !found {
		return nil, errors.New("at least one bitbucket event type must be specified")
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
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "repo:modified":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "repo:forked":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "repo:comment:added":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "repo:comment:edited":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "repo:comment:deleted":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	case "mirror:repo_synchronized":
		return nil, fmt.Errorf("'%s' not implemented", bitbucketEvent)
	default:
		return nil, fmt.Errorf("'%s' is not a valid Bitbucket Webhook event type", bitbucketEvent)
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
