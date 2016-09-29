package cli

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/mgutz/ansi"
)

const (
	stdoutColor = "green"
	stderrColor = "red"
)

func ColordChan(cmd *exec.Cmd) (stdout, stderr <-chan []byte, err error) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	//stdout = listen(outReader)
	//stderr = listen(errReader)

	stdout = listenWithColor(outReader, stdoutColor)
	stderr = listenWithColor(errReader, stderrColor)
	return
}

func listenWithColor(r io.Reader, color string) chan []byte {
	c := make(chan []byte, 0)
	go func() {
		defer close(c)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			c <- []byte(fmt.Sprintf("%s\n", ansi.Color(scanner.Text(), color)))
			//c <- []byte(fmt.Sprintf("%s", ansi.Color(scanner.Text(), color)))
		}
	}()
	return c
}
func listen(r io.Reader) chan []byte {
	c := make(chan []byte, 0)
	go func() {
		defer close(c)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			c <- scanner.Bytes()
		}
	}()
	return c
}
