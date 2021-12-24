package events

import "fmt"

// repo:refs_changed payload

type PushEvent struct {
	EventKey   `json:"eventKey"`
	EventDate  `json:"date"`
	Actor      `json:"actor"`
	Repository `json:"repository"`
	Changes    []Changes `json:"changes"`
}

func (p PushEvent) IsValid() error {
	if p.EventKey == "" {
		return fmt.Errorf("eventKey cannot be empty")
	}

	if p.EventDate == "" {
		return fmt.Errorf("date cannot be empty")
	}

	if (p.Actor == Actor{}) {
		return fmt.Errorf("actor cannot be empty")
	}

	if (p.Repository == Repository{}) {
		return fmt.Errorf("repository cannot be empty")
	}

	if len(p.Changes) == 0 || (p.Changes[0] == Changes{}) {
		return fmt.Errorf("changes cannot be empty")
	}

	return nil
}
