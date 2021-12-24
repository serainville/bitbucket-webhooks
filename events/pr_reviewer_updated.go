package events

import "fmt"

// PrReviewerUpdatedEvent maps to pr:reviewer:updated Bitbucket Webhook events
type PrReviewerUpdatedEvent struct {
	EventKey         `json:"eventKey"`
	EventDate        `json:"date"`
	Actor            `json:"actor"`
	PullRequest      `json:"pullRequest"`
	AddedReviewers   []Actor `json:"addedReviewers"`
	RemovedReviewers []Actor `json:"removedReviewers"`
}

func (p PrReviewerUpdatedEvent) Validation() error {
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

	if len(p.AddedReviewers) == 0 {
		return fmt.Errorf("addedReviewers cannot be empty")
	}

	if len(p.RemovedReviewers) == 0 {
		return fmt.Errorf("removedReviewers cannot be empty")
	}

	return nil
}
