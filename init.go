package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func initSwarm(arg *Arg) {
	writer.Add("init swarm")
	writer.Write(0, []byte(arg.String()))
	writer.Write(0, []byte("\n"))

	names := getTargetMachineNames(arg.Name)
	managers := make([]string, 0)
	workers := make([]string, 0)

	managerName := fmt.Sprintf("%s-%s", arg.Name, managerPrefix)
	workerName := fmt.Sprintf("%s-%s", arg.Name, workerPrefix)

	for _, name := range names {
		if strings.HasPrefix(name, managerName) {
			managers = append(managers, "")
		}
		if strings.HasPrefix(name, workerName) {
			workers = append(workers, "")
		}
	}
	for _, name := range names {
		if strings.HasPrefix(name, managerName) {
			index, err := strconv.Atoi(strings.Replace(name, managerName, "", 1))
			if err != nil {
				log.Fatalln(err)
			}
			managers[index] = name
		}
		if strings.HasPrefix(name, workerName) {
			index, err := strconv.Atoi(strings.Replace(name, workerName, "", 1))
			if err != nil {
				log.Fatalln(err)
			}
			workers[index] = name
		}
	}

	if len(managers) == 0 {
		log.Fatalln(errors.New("no manager exists"))
	}
	//fmt.Println(managers, workers)

	var managerAddress string
	var managerToken string
	var workerToken string
	wg := sync.WaitGroup{}
	for i, name := range managers {
		index := writer.Add(name)
		if i == 0 {
			//init manager
			managerAddress, managerToken, workerToken = initManagerNode(name, index)
			continue
		}
		wg.Add(1)
		//join manager
		go func(name string, index int) {
			joinNode(name, index, managerToken, managerAddress)
			wg.Done()
		}(name, index)
	}
	for _, name := range workers {
		index := writer.Add(name)
		wg.Add(1)
		//join worker
		go func(name string, index int) {
			joinNode(name, index, workerToken, managerAddress)
			wg.Done()
		}(name, index)
	}
	wg.Wait()
}

//execute docker-machine ssh {name} sudo docker swarm init
func initManagerNode(name string, i int) (managerAddress string, managerToken string, workerToken string) {
	args := make([]string, 0)
	args = append(args, "ssh")
	args = append(args, name)
	args = append(args, "sudo")
	args = append(args, "docker")
	args = append(args, "swarm")
	args = append(args, "init")

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	cmd := exec.Command("docker-machine", args...)
	err := execMulti(cmd, i)
	if err != nil {
		writer.Write(i, []byte(err.Error()+"\n"))
	}
	writer.Write(i, []byte("done.\n"))
	managerAddress = getManagerIP(name, i)
	managerToken = getManagerToken(name, i)
	workerToken = getWorkerToken(name, i)
	return
}

//execute docker-machine ip {name}
func getManagerIP(name string, i int) string {
	args := make([]string, 0)
	args = append(args, "ip")
	args = append(args, name)

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	out, err := exec.Command("docker-machine", args...).Output()
	if err != nil {
		log.Fatalln(err, "(", out, ")")
		return ""
	}
	res := strings.TrimSpace(string(out))
	writer.Write(i, []byte(fmt.Sprintln("manager ip:", res)))
	return res
}

//execute docker-machine ssh {name} sudo docker swarm join-token -q manager
func getManagerToken(name string, i int) string {
	args := make([]string, 0)
	args = append(args, "ssh")
	args = append(args, name)
	args = append(args, "sudo")
	args = append(args, "docker")
	args = append(args, "swarm")
	args = append(args, "join-token")
	args = append(args, "-q")
	args = append(args, "manager")

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	out, err := exec.Command("docker-machine", args...).Output()
	if err != nil {
		log.Fatalln(err, "(", out, ")")
		return ""
	}
	res := strings.TrimSpace(string(out))
	writer.Write(i, []byte(fmt.Sprintln("manager token:", res)))
	return res
}

//execute docker-machine ssh {name} sudo docker swarm join-token -q worker
func getWorkerToken(name string, i int) string {
	args := make([]string, 0)
	args = append(args, "ssh")
	args = append(args, name)
	args = append(args, "sudo")
	args = append(args, "docker")
	args = append(args, "swarm")
	args = append(args, "join-token")
	args = append(args, "-q")
	args = append(args, "worker")

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	out, err := exec.Command("docker-machine", args...).Output()
	if err != nil {
		log.Fatalln(err, "(", out, ")")
		return ""
	}
	res := strings.TrimSpace(string(out))
	writer.Write(i, []byte(fmt.Sprintln("worker token:", res)))
	return res
}

//execute docker-machine ssh {name} sudo docker swarm join --token {token} {ip}
func joinNode(name string, i int, token string, ip string) {
	args := make([]string, 0)
	args = append(args, "ssh")
	args = append(args, name)
	args = append(args, "sudo")
	args = append(args, "docker")
	args = append(args, "swarm")
	args = append(args, "join")
	args = append(args, "--token")
	args = append(args, token)
	args = append(args, ip)

	writer.Write(i, []byte(fmt.Sprintln("docker-machine", args)))
	cmd := exec.Command("docker-machine", args...)
	err := execMulti(cmd, i)
	if err != nil {
		writer.Write(i, []byte(err.Error()+"\n"))
	}
	writer.Write(i, []byte("done.\n"))
}
