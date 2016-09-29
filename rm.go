package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

func rm(arg *Arg) {
	writer.Add("remove swarm cluster")
	writer.Write(0, []byte("\n"))

	out, err := exec.Command("docker-machine", "ls", "-q", "-t", "60").Output()
	if err != nil {
		log.Println(err)
	}
	names := strings.Split(string(out), "\n")

	wg := sync.WaitGroup{}
	for _, name := range names {
		if name == "" || !strings.HasPrefix(name, arg.Name) {
			continue
		}
		index := writer.Add(name)
		wg.Add(1)
		go func(name string, index int) {
			rmNode(name, index)
			wg.Done()
		}(name, index)
	}
	wg.Wait()
}

func rmNode(name string, i int) {
	args := make([]string, 0)
	args = append(args, "rm")
	args = append(args, "-f")
	args = append(args, name)
	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	cmd := exec.Command("docker-machine", args...)
	err := execMulti(cmd, i)
	if err != nil {
		writer.Write(i, []byte(err.Error()+"\n"))
	}
	writer.Write(i, []byte("done.\n"))
}
