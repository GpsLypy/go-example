package replication

import (
	"fmt"
	"io"
)

type HeartbeatEvent struct {
	LogIdent string
}

func (e *HeartbeatEvent) Decode(data []byte) error {
	e.LogIdent = string(data)
	return nil
}

func (e *HeartbeatEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Log ident: %s\n", e.LogIdent)
	fmt.Fprintln(w)
}
