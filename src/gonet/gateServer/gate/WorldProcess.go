package gate

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/common"
	"gonet/message"
	"strconv"
)

type (
	WorldProcess struct {
		actor.Actor
		m_LostTimer *common.SimpleTimer

		m_Id uint32
	}

	IWorldlProcess interface {
		actor.IActor

		RegisterServer(int, string, int)
		SetSocketId(uint32)
	}
)

func (this *WorldProcess) RegisterServer(ServerType int, Ip string, Port int) {
	SERVER.GetWorldCluster().SendMsg(this.m_Id, "COMMON_RegisterRequest", ServerType, Ip, Port)
}

func (this *WorldProcess) SetSocketId(socketId uint32) {
	this.m_Id = socketId
}

func (this *WorldProcess) Init(num int) {
	this.Actor.Init(num)
	this.m_LostTimer = common.NewSimpleTimer(10)
	this.m_LostTimer.Start()
	this.m_Id = 0
	this.RegisterTimer(1*1000*1000*1000, this.Update)
	this.InitMessage()

	this.Actor.Start()
}

func (this *WorldProcess) Update() {
	if this.m_LostTimer.CheckTimer() {
		SERVER.GetWorldCluster().GetCluster(this.m_Id).Start()
	}
}

func (this *WorldProcess) InitMessage() {
	this.RegisterCall("COMMON_RegisterRequest", this.Handle_COMMON_RegisterRequest)

	this.RegisterCall("COMMON_RegisterResponse", this.Handle_COMMON_RegisterResponse)

	this.RegisterCall("STOP_ACTOR", this.Handle_STOP_ACTOR)

	this.RegisterCall("DISCONNECT", this.Handle_DISCONNECT)
}

func (this *WorldProcess) Handle_COMMON_RegisterRequest() {
	port, _ := strconv.Atoi(GateNetPort)
	this.RegisterServer(int(message.SERVICE_GATESERVER), GateNetIP, port)
	SERVER.GetLog().Println("客户端登录Gate服务器\n")
	fmt.Printf("客户端登录Gate服务器\n")
}

func (this *WorldProcess) Handle_COMMON_RegisterResponse() {
	//收到worldserver对自己注册的反馈
	this.m_LostTimer.Stop()
	SERVER.GetLog().Println("收到world对自己注册的反馈")
}

func (this *WorldProcess) Handle_STOP_ACTOR() {
	this.Stop()
}

func (this *WorldProcess) Handle_DISCONNECT(socketId int) {
	this.m_LostTimer.Start()
}

func DispatchPacketToClient(id int, buff []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("DispatchPacketToClient PacketFunc", err)
		}
	}()

	bitstream := base.NewBitStream(buff, len(buff))
	bitstream.ReadString() //统一格式包头名字
	//服务器标示
	if bitstream.ReadInt(base.Bit8) == int(message.SERVICE_WORLDSERVER) {
		accountId := bitstream.ReadInt64(base.Bit64)
		socketId := SERVER.GetPlayerMgr().GetSocket(accountId)
		SERVER.GetServer().SendByID(socketId, bitstream.GetBytePtr())
		return true
	}

	return false
}
