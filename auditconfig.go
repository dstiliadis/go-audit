// +build linux

package audit

import (
	"fmt"
	"strings"
)

// AuditConfig is the audit configuration
type AuditConfig struct {
	nl               *NetlinkClient
	m                *AuditMarshaller
	w                AuditWriter
	f                []AuditFilter
	SocketBufferSize int
	EventMin         uint16
	EventMax         uint16
	MaxOutOfOrder    int
	stop             chan bool
	TrackMessages    bool
	LogOutOfOrder    bool
}

func NewAuditConfig(w AuditWriter, f []AuditFilter) (*AuditConfig, error) {
	a := AuditConfig{
		w:                w,
		f:                f,
		SocketBufferSize: 16384,
		EventMin:         1300,
		EventMax:         1399,
		MaxOutOfOrder:    500,
		TrackMessages:    true,
		LogOutOfOrder:    false,
		stop:             make(chan bool),
	}

	var err error
	a.nl, err = NewNetlinkClient(a.SocketBufferSize)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize netlink %s", err)
	}

	a.m = NewAuditMarshaller(a.w, a.EventMin, a.EventMax, a.TrackMessages, a.LogOutOfOrder, a.MaxOutOfOrder, a.f)

	return &a, nil
}

// Start starts the audit process
func (a *AuditConfig) Start() {

	for {
		select {
		case <-a.stop:
			return
		default:
			msg, err := a.nl.Receive()
			if err != nil {
				continue
			}

			if msg == nil {
				continue
			}

			a.m.Consume(msg)
		}
	}
}

// Stop stops the audit process
func (a *AuditConfig) Stop() {
	a.stop <- true
}

// UpdateRules updates the rules in the system. The rules are provided
// as a list of audictl command strings
func (a *AuditConfig) UpdateRules(rules []string) error {

	if err := lExec("auditctl", "-D"); err != nil {
		return fmt.Errorf("Failed to flush existing audit rules. Error: %s", err)
	}

	if err := lExec("auditctl", "-e", "1"); err != nil {
		return fmt.Errorf("Failed to initialize. Error: %s", err)
	}

	for i, v := range rules {
		// Skip rules with no content
		if v == "" {
			continue
		}

		if err := lExec("auditctl", strings.Fields(v)...); err != nil {
			return fmt.Errorf("Failed to add rule #%d. Error: %s", i+1, err)
		}
	}

	return nil
}
