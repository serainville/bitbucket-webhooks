package events

import "fmt"

// RepModifiedEvent maps to 'repo:modified' Bitbucket Webhook events
type RepoModifiedEvent struct {
	EventKey   `json:"eventKey"`
	EventDate  `json:"date"`
	Actor      `json:"actor"`
	OldVersion RepoVersion `json:"old"`
	NewVersion RepoVersion `json:"new"`
}

func (p RepoModifiedEvent) Validation() error {
	if p.EventKey == "" {
		return fmt.Errorf("eventKey cannot be empty")
	}

	if p.EventDate == "" {
		return fmt.Errorf("eventDate cannot be empty")
	}

	if (p.Actor == Actor{}) {
		return fmt.Errorf("actor cannot be empty")
	}

	if (p.OldVersion == RepoVersion{}) {
		return fmt.Errorf("old cannot be empty")
	}

	if (p.NewVersion == RepoVersion{}) {
		return fmt.Errorf("new cannot be empty")
	}

	return nil
}
