package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/serainville/bitbucket-webhooks/signature"
)

//BitbucketEvent is an interface for all event types received from a Bitbucket Webhook
type BitbucketEvent interface {
	Validation() error
}

// EventKey stores the key for an event received by Bitbucket
type EventKey string

// EventDate stores the date an event was trigger by Bitbucket
type EventDate string

// Actor represents the actor field of a Bitbucket Webhook request
type Actor struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           uint64 `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

// PullRequest represents the pullRequest field of a Bitbucket Webhook request
type PullRequest struct {
	ID          uint64 `json:"id"`
	Version     uint64 `json:"version"`
	Title       string `json:"title"`
	State       string `json:"state"`
	Open        bool   `json:"open"`
	Closed      bool   `json:"closed"`
	CreatedDate uint64 `json:"createdDate"`
	UpdatedDate uint64 `json:"updatedDate"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
}

// Ref represents the fromRef field of a Bitbucket Webhook request
type Ref struct {
	ID           string `json:"id"`
	DisplayID    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   `json:"repository"`
}

// Repository maps to the repository key from a Bitbucket event
type Repository struct {
	Slug          string `json:"slug"`
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       `json:"project"`
	Public        bool `json:"public"`
}

// Project maps to the project key from a Bitbucket event
type Project struct {
	Key    string `json:"key"`
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
}

// Changes maps to the changes key from a Bitbucket event
type Changes struct {
	Ref struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
		Type      string `json:"type"`
	} `json:"ref"`
	RefID  string `json:"refId"`
	ToHash string `json:"toHash"`
	Type   string `json:"type"`
}

// RepoVersion maps to the version key of a Bitbucket event
type RepoVersion struct {
	Slug          string `json:"slug"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
	Project       `json:"project"`
	Public        bool `json:"public"`
}

// Participant maps to the participant key of a Bitbucket event
type Participant struct {
	Actor              `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Approved           string `json:"approved"`
	Status             string `json:"status"`
}

// PreviousTarget maps to the previousTarget key of a Bitbucket event
type PreviousTarget struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
}

// PrReviewerEvent maps to common reviewer events of a Bitbucket event
type PrReviewerEvent struct {
	EventKey       `json:"eventKey"`
	EventDate      `json:"date"`
	Actor          `json:"actor"`
	PullRequest    `json:"pullRequest"`
	Participant    `json:"participant"`
	PreviousStatus string `json:"previousStatus"`
}

// WebhookHandler is used to handle incoming Bitbuckt Webhook events
type WebhookHandler struct {
	Secret string
	// VerifySignature sets whether an HMAC signature is validated or not.
	VerifySignature bool
	signature       string
	payload         []byte
}

// WebhookEvent stores an event received from a Bitbucket Webhook request
type WebhookEvent struct {
	eventKey  string
	signature string
	payload   []byte
	Event     BitbucketEvent
}

// CreateWebookHandler with default settings
func CreateWebookHandler() *WebhookHandler {
	return &WebhookHandler{
		VerifySignature: true,
	}
}

// WebhookEvent processes an incoming Bitbucket Webhook event and returns a new WebhookEvent
func (w *WebhookHandler) WebhookEvent(resp http.ResponseWriter, req *http.Request) (WebhookEvent, error) {
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return WebhookEvent{}, fmt.Errorf("could not read body. %v", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(w.payload))

	err = w.ValidateSignature()
	if err != nil {
		return WebhookEvent{}, err
	}

	event := WebhookEvent{
		eventKey:  req.Header["X-Event-Key"][0],
		signature: req.Header["X-Hub-Signature"][0],
		payload:   payload,
	}

	if err := event.Unmarshal(); err != nil {
		return event, err
	}

	return event, nil

}

// Unmarshal a Bitbucket Webhook event into a BitbucketEvent type
// TODO: This should be refactored to deduplicate each case
func (w *WebhookEvent) Unmarshal() error {

	switch w.eventKey {
	case "diagnostic:ping":
		var event DiagnosticPingEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:opened":
		var event PullRequestEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:from_ref_updated":
		var event SourceBranchUpdatedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:modified":
		var event PrModifiedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:reviewer:updated":
		var event PrReviewerUpdatedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:reviewer:approved":
		var event PrReviewerApprovedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:reviewer:unapproved":
		var event PrReviewerUnapprovedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:reviewer:needs_work":
		var event PrReviewerNeedsWorkEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:merged":
		var event PrMergedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:declined":
		var event PrDeclinedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:deleted":
		var event PrDeletedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "pr:comment:added":
		return notImplemented()
	case "pr:comment:edited":
		return notImplemented()
	case "pr:comment:deleted":
		return notImplemented()
	case "repo:refs_changed":
		var event PushEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "repo:modified":
		var event RepoModifiedEvent
		err := json.Unmarshal(w.payload, &event)
		if err != nil {
			return err
		}
		w.Event = event
	case "repo:forked":
		return notImplemented()
	case "repo:comment:added":
		return notImplemented()
	case "repo:comment:edited":
		return notImplemented()
	case "repo:comment:deleted":
		return notImplemented()
	case "mirror:repo_synchronized":
		return notImplemented()

	default:
		return fmt.Errorf("%s is not a supported eventKey", w.eventKey)
	}

	return fmt.Errorf("failed to unmarshal event payload")
}

// ValidateSignature is used to check the authenticty of a request by checking its HMAC signature.
func (w *WebhookHandler) ValidateSignature() error {
	if len(w.Secret) > 0 && w.VerifySignature && len(w.signature) > 0 {
		return signature.Validate(w.payload, w.signature, w.Secret)
	}
	return nil
}

func notImplemented() error {
	return fmt.Errorf("not implemented")
}

// GetType returns the name of the event's struct
func GetType(event BitbucketEvent) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
