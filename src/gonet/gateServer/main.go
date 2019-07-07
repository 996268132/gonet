package main

import (
	"fmt"
	"gonet/base"
	"gonet/common"
	"gonet/gateServer/gate"
	"os"
	"os/signal"
)

func main() {
	args := os.Args
	gate.SERVER.Init()

	base.SEVERNAME = args[1]
	base.SEVERNAME = "gate"

	common.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
}
