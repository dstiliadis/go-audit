// +build linux

package audit

import (
	"encoding/json"
	"io"
	"time"
)

// AuditWriter is the interface that implements different audit writers
type AuditWriter interface {
	Write(msg *AuditMessageGroup) (err error)
}

// AuditWriterIO is an implementation of AuditWriter based on IO buffers and file system
type AuditWriterIO struct {
	e        *json.Encoder
	w        io.Writer
	attempts int
}

func NewAuditWriterIO(w io.Writer, attempts int) *AuditWriterIO {
	return &AuditWriterIO{
		e:        json.NewEncoder(w),
		w:        w,
		attempts: attempts,
	}
}

func (a *AuditWriterIO) Write(msg *AuditMessageGroup) (err error) {

	for i := 0; i < a.attempts; i++ {
		err = a.e.Encode(msg)
		if err == nil {
			break
		}

		if i != a.attempts {
			// We have to reset the encoder because write errors are kept internally and can not be retried
			a.e = json.NewEncoder(a.w)
			el.Println("Failed to write message, retrying in 1 second. Error:", err)
			time.Sleep(time.Second * 1)
		}
	}

	return err
}

// AuditWriterChannel is an channel based implementation of the interface
type AuditWriterChannel struct {
	c chan *AuditMessageGroup
}

// NewAuditWriterChannel creates a new audit writer with a provided channel
func NewAuditWriterChannel(c chan *AuditMessageGroup) *AuditWriterChannel {
	return &AuditWriterChannel{
		c: c,
	}
}

// Write is the implementation of the Write function of the interface
func (a *AuditWriterChannel) Write(msg *AuditMessageGroup) (err error) {
	a.c <- msg

	return nil
}
