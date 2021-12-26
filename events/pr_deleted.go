package events

import "fmt"

// PrDeletedEvent maps to 'pr:deleted' Bitbucket Webhook events
type PrDeletedEvent struct {
	EventKey    `json:"eventKey"`
	EventDate   `json:"date"`
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

// Validation checks whether a pr:deleted event is valid
func (p PrDeletedEvent) Validation() error {
	if p.EventKey == "" {
		return fmt.Errorf("eventKey cannot be empty")
	}

	if p.EventDate == "" {
		return fmt.Errorf("eventDate cannot be empty")
	}

	if (p.Actor == Actor{}) {
		return fmt.Errorf("actor cannot be empty")
	}

	if (p.PullRequest == PullRequest{}) {
		return fmt.Errorf("pullRequest cannot be empty")
	}

	return nil
}
