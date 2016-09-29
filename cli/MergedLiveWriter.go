package cli

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

type MergedLiveWriter struct {
	bufs           []*bytes.Buffer
	names          []string
	lineBreakCount int
	sep            string
	mu             sync.Mutex
}

func NewMergedLiveWriter() *MergedLiveWriter {
	res := &MergedLiveWriter{
		bufs:  make([]*bytes.Buffer, 0),
		names: make([]string, 0),
		sep:   "*******************************************************************************",
	}
	return res
}

func (this *MergedLiveWriter) Add(name string) int {
	this.bufs = append(this.bufs, bytes.NewBufferString(""))
	this.names = append(this.names, name)
	return len(this.bufs) - 1
}
func (this *MergedLiveWriter) writeHeader(index int) {
	buf := this.bufs[index]
	name := this.names[index]
	sep := this.sep[0 : len(this.sep)-1-len(name)-1]
	buf.WriteString(name)
	buf.WriteString(" ")
	buf.WriteString(sep)
	buf.WriteString("\n")
}

func (this *MergedLiveWriter) String() string {
	all := bytes.NewBufferString("")
	for _, buf := range this.bufs {
		all.Write(buf.Bytes())
		all.WriteString("\n")
	}
	return all.String()
}

func (this *MergedLiveWriter) Write(index int, p []byte) error {
	defer this.mu.Unlock()
	this.mu.Lock()

	//clear text by lineBreak
	//if this.lineBreakCount > 0 {
	//	fmt.Printf("\033[%dF", this.lineBreakCount)
	//	fmt.Print("\033[J")
	//}
	//clear all text
	//fmt.Print("\033[10000S")
	//fmt.Print("\033[0,0H")
	fmt.Print("\033[3J")
	fmt.Print("\033c")

	buf := this.bufs[index]
	if buf.Len() == 0 {
		this.writeHeader(index)
		this.lineBreakCount++
	}
	this.lineBreakCount += strings.Count(string(p), "\n") + len(this.bufs)

	_, err := buf.Write(p)
	if err != nil {
		return err
	}
	fmt.Print(this.String())
	return nil
}
