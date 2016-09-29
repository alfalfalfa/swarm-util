package main

import (
	"os/exec"

	"log"
	"strings"

	"github.com/alfalfalfa/swarm-util/cli"
	"github.com/alfalfalfa/swarm-util/util"
)

const (
	managerPrefix = "m"
	workerPrefix  = "w"
)

var writer *cli.MergedLiveWriter

func execMulti(cmd *exec.Cmd, index int) (string, error) {
	stdout, stderr, err := cli.ColordChan(cmd)
	if err != nil {
		return "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	out := util.MargeChan(stdout, stderr)
	go func() {
		for {
			b := <-out
			if b == nil {
				return
			}
			writer.Write(index, b)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return "", err
	}
	return writer.StringByIndex(index), nil
}

func getTargetMachineNames(targetName string) []string {
	res := make([]string, 0)
	out, err := exec.Command("docker-machine", "ls", "-q", "-t", "60").Output()
	if err != nil {
		log.Fatalln(err, "(", out, ")")
		return res
	}
	names := strings.Split(string(out), "\n")

	for _, name := range names {
		if name == "" || !strings.HasPrefix(name, targetName) {
			continue
		}
		res = append(res, name)
	}

	return res
}
