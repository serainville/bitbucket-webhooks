package events

import "fmt"

// SourceBranchUpdatedEvent maps to a from_ref_updated Bitbucket Webhook event
type SourceBranchUpdatedEvent struct {
	EventKey         `json:"eventKey"`
	EventDate        `json:"date"`
	Actor            `json:"actor"`
	PullRequest      `json:"pullRequest"`
	PreviousFromHash string `json:"previousFramHash"`
}

func (p SourceBranchUpdatedEvent) IsValid() error {
	if p.EventKey == "" {
		return fmt.Errorf("eventKey cannot be empty")
	}

	if p.EventDate == "" {
		return fmt.Errorf("date cannot be empty")
	}

	if (p.Actor == Actor{}) {
		return fmt.Errorf("actor cannot be empty")
	}

	if (p.PullRequest == PullRequest{}) {
		return fmt.Errorf("pullRequest cannot be empty")
	}

	if p.PreviousFromHash == "" {
		return fmt.Errorf("previousFromHash cannot be empty")
	}

	return nil
}
