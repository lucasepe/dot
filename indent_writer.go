package dot

import (
	"fmt"
	"io"
)

const (
	space = "  "
)

// IndentWriter is a writer with indentation.
type IndentWriter struct {
	level  int
	writer io.Writer
}

// NewIndentWriter returns a IndentWriter from a io.Writer.
func NewIndentWriter(w io.Writer) *IndentWriter {
	return &IndentWriter{level: 0, writer: w}
}

// Indent writes a tab `\t` to the writer
func (i *IndentWriter) Indent() {
	i.level++
	fmt.Fprint(i.writer, space)
}

// BackIndent decrements the current indentation level.
func (i *IndentWriter) BackIndent() {
	i.level--
}

// IndentWhile writes an indented block.
func (i *IndentWriter) IndentWhile(block func()) {
	i.Indent()
	block()
	i.BackIndent()
}

// NewLineIndentWhile writes an indented block between new line chars `\n`.
func (i *IndentWriter) NewLineIndentWhile(block func()) {
	i.NewLine()
	i.Indent()
	block()
	i.BackIndent()
	i.NewLine()
}

// NewLine add an indented new line.
func (i *IndentWriter) NewLine() {
	fmt.Fprint(i.writer, "\n")
	for j := 0; j < i.level; j++ {
		fmt.Fprint(i.writer, space)
	}
}

// Write makes it an io.Writer
func (i *IndentWriter) Write(data []byte) (n int, err error) {
	return i.writer.Write(data)
}

// WriteString writes an indented string
func (i *IndentWriter) WriteString(s string) (n int, err error) {
	fmt.Fprint(i.writer, s)
	return len(s), nil
}
