package main

import (
	"fmt"
	"gonet/actor"
	"gonet/base"
	"gonet/message"
	"gonet/network"

	"github.com/golang/protobuf/proto"
)

type (
	EventProcess struct {
		actor.Actor

		Client      *network.ClientSocket
		AccountId   int64
		PlayerId    int64
		AccountName string
		SimId       int64
	}

	IEventProcess interface {
		actor.IActor
		LoginGame()
		LoginAccount()
		SendPacket(proto.Message)
	}
)

func SendPacket(packet proto.Message) {
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	CLIENT.Send(buff)
}

func (this *EventProcess) SendPacket(packet proto.Message) {
	buff := message.Encode(packet)
	buff = base.SetTcpEnd(buff)
	this.Client.Send(buff)
}

func (this *EventProcess) PacketFunc(socketid int, buff []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("EventProcess PacketFunc", err)
		}
	}()

	packetId, data := message.Decode(buff)
	packet := message.GetPakcet(packetId)
	if packet == nil {
		return true
	}
	err := message.UnmarshalText(packet, data)
	if err == nil {
		bitstream := base.NewBitStream(make([]byte, 1024), 1024)
		if !message.GetMessagePacket(packet, bitstream) {
			return true
		}
		var io actor.CallIO
		io.Buff = bitstream.GetBuffer()
		io.SocketId = socketid
		this.Send(io)
		return true
	}

	return true
}

func (this *EventProcess) Init(num int) {
	this.Actor.Init(num)
	this.RegisterCall("W_C_SelectPlayerResponse", func(packet *message.W_C_SelectPlayerResponse) {
		fmt.Println("EventProcess PacketFunc: ", "W_C_SelectPlayerResponse")
		this.AccountId = packet.GetAccountId()
		nLen := len(packet.GetPlayerData())
		//fmt.Println(len(packet.PlayerData), this.AccountId, packet.PlayerData)
		if nLen == 0 {
			packet1 := &message.C_W_CreatePlayerRequest{PacketHead: message.BuildPacketHead(this.AccountId, int(message.SERVICE_WORLDSERVER)),
				PlayerName: proto.String("我是大坏蛋"),
				Sex:        proto.Int32(int32(0))}
			this.SendPacket(packet1)
		} else {
			this.PlayerId = packet.GetPlayerData()[0].GetPlayerID()
			this.LoginGame()
		}
	})

	this.RegisterCall("W_C_CreatePlayerResponse", func(packet *message.W_C_CreatePlayerResponse) {
		if packet.GetError() == 0 {
			this.PlayerId = packet.GetPlayerId()
			this.LoginGame()
		} else { //创建失败

		}
	})

	this.RegisterCall("A_C_LoginRequest", func(packet *message.A_C_LoginRequest) {
		if packet.GetError() == base.ACCOUNT_NOEXIST {
			packet1 := &message.C_A_RegisterRequest{PacketHead: message.BuildPacketHead(0, int(message.SERVICE_ACCOUNTSERVER)),
				AccountName: packet.AccountName, SocketId: proto.Int32(0)}
			this.SendPacket(packet1)
		}
	})

	this.RegisterCall("A_C_RegisterResponse", func(packet *message.A_C_RegisterResponse) {
		//注册失败
		if packet.GetError() != 0 {
		}
	})

	this.RegisterCall("W_C_ChatMessage", func(packet *message.W_C_ChatMessage) {
		fmt.Println("收到【", packet.GetSenderName(), "】发送的消息[", packet.GetMessage()+"]")
	})

	this.Actor.Start()
}

func (this *EventProcess) LoginGame() {
	packet1 := &message.C_W_Game_LoginRequset{PacketHead: message.BuildPacketHead(this.AccountId, int(message.SERVICE_WORLDSERVER)),
		PlayerId: proto.Int64(this.PlayerId)}
	this.SendPacket(packet1)
}

var (
	id int
)

func (this *EventProcess) LoginAccount() {
	id++
	this.AccountName = fmt.Sprintf("test%d", id)
	//this.AccountName = fmt.Sprintf("test%d", base.RAND().RandI(0, 7000))
	packet1 := &message.C_A_LoginRequest{PacketHead: message.BuildPacketHead(0, int(message.SERVICE_ACCOUNTSERVER)),
		AccountName: proto.String(this.AccountName), BuildNo: proto.String(base.BUILD_NO), SocketId: proto.Int32(0)}
	this.SendPacket(packet1)
}

var (
	PACKET *EventProcess
)
