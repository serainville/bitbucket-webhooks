package events

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type BitbucketEvent interface {
	Validation() error
}

type EventKey string
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
	FromRef     `json:"fromRef"`
	ToRef       `json:"toRef"`
}

// FromRef represents the fromRef field of a Bitbucket Webhook request
type FromRef struct {
	ID           string `json:"id"`
	DisplayId    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   `json:"repository"`
}

type ToRef struct {
	ID           string `json:"id"`
	DisplayId    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   `json:"repository"`
}

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

type Project struct {
	Key    string `json:"key"`
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	Type   string `json:"type"`
}

type Changes struct {
	Ref struct {
		ID        string `json:"id"`
		DisplayID string `json:"displayId"`
		Type      string `json:"type"`
	} `json:"ref"`
	RefId  string `json:"refId"`
	ToHash string `json:"toHash"`
	Type   string `json:"type"`
}

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

type Participant struct {
	Actor              `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Approved           string `json:"approved"`
	Status             string `json:"status"`
}

type PreviousTarget struct {
	ID              string `json:"id"`
	DisplayId       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
}

type PrReviewerEvent struct {
	EventKey       `json:"eventKey"`
	EventDate      `json:"date"`
	Actor          `json:"actor"`
	PullRequest    `json:"pullRequest"`
	Participant    `json:"participant"`
	PreviousStatus string `json:"previousStatus"`
}

func NewBitbucketEvent(eventKey string, payload []byte) (BitbucketEvent, error) {
	switch eventKey {
	case "diagnostic:ping":
		var event DiagnosticPingEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:opened":
		var event PullRequestEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:from_ref_updated":
		var event SourceBranchUpdatedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:modified":
		var event PrModifiedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:reviewer:updated":
		var event PrReviewerUpdatedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:reviewer:approved":
		var event PrReviewerApprovedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:reviewer:unapproved":
		var event PrReviewerUnapprovedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:reviewer:needs_work":
		var event PrReviewerNeedsWorkEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:merged":
		var event PrMergedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:declined":
		var event PrDeclinedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:deleted":
		var event PrDeletedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "pr:comment:added":
		return nil, notImplemented()
	case "pr:comment:edited":
		return nil, notImplemented()
	case "pr:comment:deleted":
		return nil, notImplemented()
	case "repo:refs_changed":
		var event PushEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "repo:modified":
		var event RepoModifiedEvent
		err := json.Unmarshal(payload, &event)
		return event, err
	case "repo:forked":
		return nil, notImplemented()
	case "repo:comment:added":
		return nil, notImplemented()
	case "repo:comment:edited":
		return nil, notImplemented()
	case "repo:comment:deleted":
		return nil, notImplemented()
	case "mirror:repo_synchronized":
		return nil, notImplemented()

	default:
		return nil, fmt.Errorf("%s is not a supported eventKey", eventKey)
	}
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
