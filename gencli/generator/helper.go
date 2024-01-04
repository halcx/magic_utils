package generator

import (
	"fmt"
	"io"
)

type codeWriter struct {
	Out io.Writer
}

// Line writes a single line.
func (c *codeWriter) Line(line string) {
	_, _ = fmt.Fprintln(c.Out, line)
}

// Linef writes a single line with formatting (as per fmt.Sprintf).
func (c *codeWriter) Linef(line string, args ...interface{}) {
	_, _ = fmt.Fprintf(c.Out, line+"\n", args...)
}
