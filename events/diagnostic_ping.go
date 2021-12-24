package events

import (
	"fmt"
)

type DiagnosticPingEvent struct {
	Test bool `json:"test"`
}

func (dp DiagnosticPingEvent) Validation() error {
	if !dp.Test {
		return fmt.Errorf("test must be true")
	}
	return nil
}