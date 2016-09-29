package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type LiveWriter struct {
	buf            *bytes.Buffer
	lineBreakCount int
}

func NewLiveWriter() io.Writer {
	res := &LiveWriter{
		buf: bytes.NewBufferString(""),
	}
	return res
}

func (this *LiveWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	//fmt.Print(strings.Repeat("\r\r", this.buf.Len()))
	//fmt.Print("\033[2J")
	//fmt.Print("\033[2K")
	//fmt.Printf("\033[%dP", this.buf.Len())
	if this.lineBreakCount > 0 {
		fmt.Printf("\033[%dF", this.lineBreakCount)
		fmt.Print("\033[J")
	}
	this.lineBreakCount += strings.Count(string(p), "\n")

	n, err = this.buf.Write(p)
	if err != nil {
		return
	}
	fmt.Print(this.buf.String())
	return
}
