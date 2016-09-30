package main

import (
	"os/exec"

	"log"
	"strings"

	"regexp"

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

func findIpAddressAndPort(str string) string {
	r := regexp.MustCompile(`(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+`)
	for _, s := range strings.Split(str, "\n") {
		if !strings.Contains(s, ".") || !strings.Contains(s, ":") {
			continue
		}
		matches := r.FindAllStringSubmatch(s, -1)
		if len(matches) == 0 {
			continue
		}
		if len(matches[0]) == 0 {
			continue
		}
		return matches[0][0]
	}
	return ""
}
