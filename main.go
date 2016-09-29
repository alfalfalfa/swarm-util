package main

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/alfalfalfa/goarg"
)

func main() {
	arg := &Arg{}
	err := goarg.Parse(arg)
	if err != nil {
		log.Fatalln(err)
	}

	b, e := json.Marshal(arg)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(string(b))
}

type Arg struct {
	Int    int     `arg:"0" usage:"this is int"`
	String string  `arg:"1" usage:"this is string"`
	Float  float64 `arg:"2" usage:"this is float"`
	Bool   bool    `name:"c" usage:"this is bool"`
}
