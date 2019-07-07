package main

import (
	"gonet/actor"
	"gonet/common"
	"gonet/worldServer/world"
	"gonet/worldServer/world/chat"
	"gonet/worldServer/world/cmd"
	"gonet/worldServer/world/data"
	"gonet/worldServer/world/mail"
	"gonet/worldServer/world/player"
	"gonet/worldServer/world/social"
	"gonet/worldServer/world/toprank"
)

func InitMgr(serverName string) {
	//一些共有数据量初始化
	common.Init()
	cmd.Init()
	data.InitRepository()
	player.PLAYERMGR.Init(1000)
	chat.CHATMGR.Init(1000)
	mail.MAILMGR.Init(1000)
	toprank.MGR().Init(1000)
	player.PLAYERSIMPLEMGR.Init(1000)
	social.MGR().Init(1000)
	actor.MGR.InitActorHandle(world.SERVER.GetServer())
}
