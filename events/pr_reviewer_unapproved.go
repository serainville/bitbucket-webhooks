package events

import "fmt"

// PrReviewerUnapprovedEvent maps to 'pr:reviewer:unapproved' Bitbucket Webhook events
type PrReviewerUnapprovedEvent PrReviewerEvent

func (p PrReviewerUnapprovedEvent) Validation() error {
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

	if (p.Participant == Participant{}) {
		return fmt.Errorf("participant cannot be empty")
	}

	if len(p.PreviousStatus) == 0 {
		return fmt.Errorf("previousStatus cannot be empty")
	}

	return nil
}
