package events

import "fmt"

type ReviewerApprovedEvent PrReviewerEvent

func (p ReviewerApprovedEvent) IsValid() error {
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
