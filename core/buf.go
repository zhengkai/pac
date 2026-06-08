package core

import (
	"bytes"
	"fmt"
)

type Buf struct {
	buf bytes.Buffer
}

func (b *Buf) w(s string) {
	b.buf.WriteString(s)
	b.br()
}

func (b *Buf) wf(format string, args ...interface{}) {
	fmt.Fprintf(&b.buf, format, args...)
	b.br()
}

func (b *Buf) br() {
	b.buf.Write([]byte("\n"))
}
