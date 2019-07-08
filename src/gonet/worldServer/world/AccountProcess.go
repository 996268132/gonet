package world

import (
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/message"
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_Id uint32
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(uint32)
	}
)

func (this *AccountProcess) SetSocketId(socketId uint32) {
	this.m_Id = socketId
}

func (this *AccountProcess) RegisterServer(ServerType int, Ip string, Port int) {
	SERVER.GetAccountCluster().SendMsg(this.m_Id, "COMMON_RegisterRequest", ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.RegisterTimer(1*1000*1000*1000, this.Update)
	this.InitMessage()

	this.Actor.Start()
}

func (this *AccountProcess) InitMessage() {
	this.RegisterCall("COMMON_RegisterRequest", this.Handle_COMMON_RegisterRequest)

	this.RegisterCall("COMMON_RegisterResponse", this.Handle_COMMON_RegisterResponse)

	this.RegisterCall("DISCONNECT", this.Handle_DISCONNECT)

	this.RegisterCall("STOP_ACTOR", this.Handle_STOP_ACTOR)

	this.RegisterCall("G_ClientLost", this.Handle_G_ClientLost)

	this.RegisterCall("A_W_CreatePlayer", this.Handle_A_W_CreatePlayer)
}

func (this *AccountProcess) Update() {
	if this.m_LostTimer.CheckTimer() {
		SERVER.GetAccountCluster().GetCluster(this.m_Id).Start()
	}
}

func (this *AccountProcess) Handle_COMMON_RegisterRequest() {
	this.RegisterServer(int(message.SERVICE_WORLDSERVER), WorldNetIP, base.Int(WorldNetPort))
}

func (this *AccountProcess) Handle_COMMON_RegisterResponse() {
	this.m_LostTimer.Stop()
}

func (this *AccountProcess) Handle_DISCONNECT(socketId int) {
	this.m_LostTimer.Start()
}

func (this *AccountProcess) Handle_STOP_ACTOR() {
	this.Stop()
}

func (this *AccountProcess) Handle_G_ClientLost(accountId int64) {
	SERVER.GetServer().CallMsg("G_ClientLost", accountId)
}

func (this *AccountProcess) Handle_A_W_CreatePlayer(accountId int64, playerId int64, playername string, sex int32, socketId int) {
	SERVER.GetServer().CallMsg("A_W_CreatePlayer", accountId, playerId, playername, sex, socketId)
}
