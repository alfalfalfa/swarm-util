package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func leaveSwarm(arg *Arg) {
	writer.Add("leave swarm")
	writer.Write(0, []byte(arg.String()))
	writer.Write(0, []byte("\n"))

	wg := sync.WaitGroup{}
	for _, name := range getTargetMachineNames(arg.Name) {
		index := writer.Add(name)
		wg.Add(1)
		go func(name string, index int) {
			leaveNode(name, index)
			wg.Done()
		}(name, index)
	}
	wg.Wait()
}

//execute docker-machine ssh {name} sudo docker swarm leave
func leaveNode(name string, i int) {
	args := make([]string, 0)
	args = append(args, "ssh")
	args = append(args, name)
	args = append(args, "sudo")
	args = append(args, "docker")
	args = append(args, "swarm")
	args = append(args, "leave")
	args = append(args, "--force")

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	cmd := exec.Command("docker-machine", args...)
	_, err := execMulti(cmd, i)
	if err != nil {
		writer.Write(i, []byte(err.Error()+"\n"))
	}
	writer.Write(i, []byte("done.\n"))
}
