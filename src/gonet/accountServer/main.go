package main

import (
	"fmt"
	"gonet/accountServer/account"
	"gonet/base"
	"gonet/common"
	"os"
	"os/signal"
)

func main() {
	args := os.Args
	account.SERVER.Init()

	base.SEVERNAME = args[1]
	base.SEVERNAME = "account"
	common.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Printf("server【%s】 exit ------- signal:[%v]", args[1], s)
}
