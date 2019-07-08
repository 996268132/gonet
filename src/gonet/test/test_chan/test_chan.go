package main

import (
	"fmt"
	"math/rand"
	"time"
)

type (
	TPlayer struct {
		name  string
		score int
	}

	GateMgr struct {
		loginChan chan TPlayer
	}
)

func NewGateMgr() *GateMgr {
	fmt.Printf("Gate Init...\n")
	gate := &GateMgr{}
	gate.loginChan = make(chan TPlayer, 10)
	return gate
}

func (this *GateMgr) AddPlayer(p TPlayer) {
	fmt.Printf("玩家：%s 准备登录...\n", p.name)
	this.loginChan <- p
}

func (this *GateMgr) Logining() {
	for {

		player := TPlayer{}
		for i := 0; i < 1000; i++ {
			player.name = fmt.Sprintf("test_%d", i)
			player.score = rand.Int() % 1000
			this.AddPlayer(player)
		}

		time.Sleep(30000)

	}
}

func (this *GateMgr) ProcessLogin(index int) {
	fmt.Printf("开始处理登录...\n")
	for {
		select {
		case player, ok := <-this.loginChan:

			if ok == false {
				fmt.Println("chan is close")
				break
			} else {
				fmt.Printf("线程 %d, 玩家：%s 已登录，玩家积分：%d\n", index, player.name, player.score)
			}
		}
	}
}

func (this *GateMgr) Init() {

	go this.Logining()
	for i := 0; i < 100; i++ {
		go this.ProcessLogin(i)
	}

}

func main() {
	fmt.Printf("test init\n")

	Gate_ := NewGateMgr()
	Gate_.Init()

	for {
		time.Sleep(2 * time.Second)
	}
}
