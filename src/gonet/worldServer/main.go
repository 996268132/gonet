package main

import (
	"fmt"
	"gonet/base"
	"gonet/worldServer/world"
	"os"
	"os/signal"
)

func main() {
	args := os.Args
	world.SERVER.Init()

	base.SEVERNAME = args[1]
	base.SEVERNAME = "world"
	InitMgr(args[1])

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
}
