package main

import (
	"fmt"
	"os/exec"
	"sync"
)

func create(arg *Arg) {
	writer.Add("create swarm cluster")
	writer.Write(0, []byte(arg.String()))
	writer.Write(0, []byte("\n"))
	//fmt.Println("create swarm cluster")
	//fmt.Println(arg.String())

	//create nodes
	wg := sync.WaitGroup{}
	createNodes(&wg, arg, managerPrefix, arg.Manager)
	createNodes(&wg, arg, workerPrefix, arg.Worker)
	wg.Wait()
}

func createNodes(wg *sync.WaitGroup, arg *Arg, prefix string, count int) {
	for i := 0; i < count; i++ {
		nodeName := fmt.Sprintf("%s-%s%d", arg.Name, prefix, i)
		index := writer.Add(nodeName)
		wg.Add(1)
		go func(name string, index int) {
			createNode(name, arg.Options, index)
			wg.Done()
		}(nodeName, index)
	}
}

//execute docker-machine create
func createNode(name string, options []string, i int) {
	args := make([]string, 0)
	args = append(args, "create")
	args = append(args, options...)
	args = append(args, name)
	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	//writer.Write(i, []byte(fmt.Sprintln("docker-machine")))

	cmd := exec.Command("docker-machine", args...)
	//cmd := exec.Command("test.sh")
	_, err := execMulti(cmd, i)
	if err != nil {
		writer.Write(i, []byte(err.Error()+"\n"))
	}
	writer.Write(i, []byte("done.\n"))
}
