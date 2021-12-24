package events

import "fmt"

// pr:modified payload

type PrModifiedEvent struct {
	EventKey            `json:"eventKey"`
	EventDate           `json:"date"`
	Actor               `json:"actor"`
	PullRequest         `json:"pullRequest"`
	PreviousTitle       string `json:"previousTitle"`
	PreviousDescription string `json:"previousDescription"`
	PreviousTarget      `json:"previousTarget"`
}

func (p PrModifiedEvent) IsValid() error {
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
