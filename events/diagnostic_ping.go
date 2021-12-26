package events

import (
	"fmt"
)

// DiagnosticPingEvent maps to a Bitbucket diagnostic ping event, typically used for testing a webhook
// end point.
type DiagnosticPingEvent struct {
	Test bool `json:"test"`
}

// Validation checks with a DiagnosticPingEvent is valid
func (dp DiagnosticPingEvent) Validation() error {
	if !dp.Test {
		return fmt.Errorf("test must be true")
	}
	return nil
}
