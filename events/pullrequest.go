package events

import (
	"fmt"
)

type PullRequestEvent struct {
	EventKey    `json:"eventKey"`
	EventDate   `json:"date"`
	Actor       `json:"actor"`
	PullRequest `json:"pullRequest"`
}

func (p PullRequestEvent) IsValid() error {
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
