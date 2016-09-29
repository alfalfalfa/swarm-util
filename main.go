package main

import (
	"log"

	"flag"

	"encoding/json"

	"fmt"
	"os/exec"
	"sync"

	"github.com/alfalfalfa/goarg"
	"github.com/alfalfalfa/swarm-util/cli"
	"github.com/alfalfalfa/swarm-util/util"
)

func main2() {
	wg := sync.WaitGroup{}
	writer := cli.NewMergedLiveWriter()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		writer.Add(fmt.Sprint(i))
		go func(i int) {
			defer wg.Done()

			cmd := exec.Command("test.sh")
			stdout, stderr, err := cli.ColordChan(cmd)
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			out := util.MargeChan(stdout, stderr)

			go func() {
				for {
					b := <-out
					if b == nil {
						return
					}
					writer.Write(i, b)
				}
			}()
			err = cmd.Wait()
			if err != nil {
				log.Fatal(err)
			}
		}(i)
	}
	wg.Wait()
}
func main() {
	writer = cli.NewMergedLiveWriter()

	arg := &Arg{}
	err := goarg.Parse(arg)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(arg.String())
	if arg.Command == "" {
		flag.Usage()
		return
	}

	switch arg.Command {
	case "create":
		if arg.Name == "" {
			flag.Usage()
			return
		}
		if arg.Manager == 0 {
			flag.Usage()
			return
		}
		create(arg)
	case "rm":
		if arg.Name == "" {
			flag.Usage()
			return
		}
		rm(arg)
	case "init":
		if arg.Name == "" {
			flag.Usage()
			return
		}
		initSwarm(arg)
	}
}

type Arg struct {
	Command string   `arg:"0" usage:"sub command:[create, ]"`
	Options []string `arg:"*" usage:"sub command arguments"`
	Name    string   `name:"name" usage:"name of swarm cluster"`

	//create
	Manager int `name:"manager" usage:"number of swarm manager node"`
	Worker  int `name:"worker" usage:"number of swarm worker node"`
}

func (this *Arg) String() string {
	b, e := json.Marshal(this)
	if e != nil {
		log.Fatalln(e)
	}
	return string(b)
}
