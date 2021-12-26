package events

import "fmt"

// PrModifiedEvent maps to 'pr:modified' Bitbucket Webhook events
type PrModifiedEvent struct {
	EventKey            `json:"eventKey"`
	EventDate           `json:"date"`
	Actor               `json:"actor"`
	PullRequest         `json:"pullRequest"`
	PreviousTitle       string `json:"previousTitle"`
	PreviousDescription string `json:"previousDescription"`
	PreviousTarget      `json:"previousTarget"`
}

// Validation checks whether a pr:modified event is valid
func (p PrModifiedEvent) Validation() error {
	if p.EventKey == "" {
		return fmt.Errorf("eventKey cannot be empty")
	}

	if p.EventDate == "" {
		return fmt.Errorf("eventDate cannot be empty")
	}

	if (p.Actor == Actor{}) {
		return fmt.Errorf("actor cannot be empty")
	}

	if p.PreviousTitle == "" {
		return fmt.Errorf("previousTitle cannot be empty")
	}

	if p.PreviousDescription == "" {
		return fmt.Errorf("previousDescription cannot be empty")
	}

	if (p.PullRequest == PullRequest{}) {
		return fmt.Errorf("pullRequest cannot be empty")
	}

	if (p.PreviousTarget == PreviousTarget{}) {
		return fmt.Errorf("previousTarget cannot be empty")
	}

	return nil
}
