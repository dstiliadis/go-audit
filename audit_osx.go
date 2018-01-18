// +build darwin !linux

package audit

import (
	"time"
)

// AuditMarshaller is Mac mock
type AuditMarshaller struct {
}

// AuditFilter is an OSX mock
type AuditFilter struct {
}

// AuditMessageGroup is an OSX mock
type AuditMessageGroup struct {
	Seq           int               `json:"sequence"`
	AuditTime     string            `json:"timestamp"`
	CompleteAfter time.Time         `json:"-"`
	Msgs          []*AuditMessage   `json:"messages"`
	UidMap        map[string]string `json:"uid_map"`
	Syscall       string            `json:"-"`
}

// AuditMessage is an OSX mock
type AuditMessage struct {
	Type      uint16 `json:"type"`
	Data      string `json:"data"`
	Seq       int    `json:"-"`
	AuditTime string `json:"-"`
}

// AuditWriter is an OSX mock
type AuditWriter interface {
	Write(msg *AuditMessageGroup) (err error)
}

// AuditWriterChannel is an channel based implementation of the interface
type AuditWriterChannel struct {
	c chan *AuditMessageGroup
}

// NetLinkClient is an OSX mock
type NetlinkClient struct{}

// NewNetlinkClient creates a new netlink client
func NewNetlinkClient(recvSize int) (*NetlinkClient, error) {
	return nil, nil
}

// NewAuditMarshaller creates a new marshaller
func NewAuditMarshaller(w AuditWriter, eventMin uint16, eventMax uint16, trackMessages, logOOO bool, maxOOO int, filters []AuditFilter) *AuditMarshaller {
	return nil
}

// NewAuditWriterChannel is the channel implementation of the writer
func NewAuditWriterChannel(c chan *AuditMessageGroup) *AuditWriterChannel {
	return nil
}

// Write is the implementation of the Write function of the interface
func (a *AuditWriterChannel) Write(msg *AuditMessageGroup) (err error) {
	return nil
}

// Run runs the loop
func Run(nl *NetlinkClient, m *AuditMarshaller) {

}
