package gate

import (
	"gonet/actor"
	"gonet/common"
	"gonet/message"
	"strconv"
)

type (
	AccountProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_Id uint32
	}

	IAccountProcess interface {
		actor.IActor

		RegisterServer(int, int, string, int)
		SetSocketId(uint32)
	}
)

func (this *AccountProcess) SetSocketId(socketId uint32) {
	this.m_Id = socketId
}

func (this *AccountProcess) RegisterServer(ServerType int, Ip string, Port int) {
	SERVER.GetAccountCluster().GetCluster(this.m_Id).SendMsg("COMMON_RegisterRequest", ServerType, Ip, Port)
}

func (this *AccountProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.RegisterTimer(1*1000*1000*1000, this.Update)

	this.InitMessage()

	this.Actor.Start()
}

func (this *AccountProcess) Update() {
	if this.m_LostTimer.CheckTimer() {
		SERVER.GetAccountCluster().GetCluster(this.m_Id).Start()
	}
}

func (this *AccountProcess) InitMessage() {
	this.RegisterCall("COMMON_RegisterRequest", this.Handle_COMMON_RegisterRequest)

	this.RegisterCall("COMMON_RegisterResponse", this.Handle_COMMON_RegisterResponse)

	this.RegisterCall("STOP_ACTOR", this.Handle_STOP_ACTOR)

	this.RegisterCall("DISCONNECT", this.Handle_DISCONNECT)

	this.RegisterCall("A_G_Account_Login", this.Handle_A_G_Account_Login)

	this.RegisterCall("A_C_RegisterResponse", this.Handle_A_C_RegisterResponse)

	this.RegisterCall("A_C_LoginRequest", this.Handle_A_C_LoginRequest)
}

func (this *AccountProcess) Handle_COMMON_RegisterRequest() {
	port, _ := strconv.Atoi(GateNetPort)
	this.RegisterServer(int(message.SERVICE_GATESERVER), GateNetIP, port)
}

func (this *AccountProcess) Handle_COMMON_RegisterResponse() {
	this.m_LostTimer.Stop()
	SERVER.GetPlayerMgr().SendMsg("Account_Relink")
}

func (this *AccountProcess) Handle_STOP_ACTOR() {
	this.Stop()
}

func (this *AccountProcess) Handle_DISCONNECT(socketId int) {
	this.m_LostTimer.Start()
}

func (this *AccountProcess) Handle_A_G_Account_Login(accountId int64, socketId int) {
	SERVER.GetPlayerMgr().SendMsg("ADD_ACCOUNT", accountId, socketId)
}

func (this *AccountProcess) Handle_A_C_RegisterResponse(packet *message.A_C_RegisterResponse) {
	buff := message.Encode(packet)
	SERVER.GetServer().SendByID(int(packet.GetSocketId()), buff)
}

func (this *AccountProcess) Handle_A_C_LoginRequest(packet *message.A_C_LoginRequest) {
	buff := message.Encode(packet)
	SERVER.GetServer().SendByID(int(packet.GetSocketId()), buff)
}
