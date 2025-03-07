package player

import (
	"database/sql"
	"fmt"
	"gonet/actor"
	"gonet/common"
	"gonet/db"
	"gonet/message"
	"gonet/worldServer/world"

	"github.com/golang/protobuf/proto"
)

type (
	Player struct {
		actor.Actor

		PlayerData
		m_db           *sql.DB
		m_offlineTimer *common.SimpleTimer
	}

	IPlayer interface {
		actor.IActor

		Update()
		//IsOffline() bool
		//IsLogout() bool
		//IsInGame() bool
	}
)

func (this *Player) Init(num int) {
	this.Actor.Init(num)
	this.PlayerData.Init()
	this.RegisterTimer(1000*1000*1000, this.Update) //定时器
	this.m_offlineTimer = common.NewSimpleTimer(5 * 60)
	this.m_db = world.SERVER.GetDB()
	this.InitMessage()
	this.Actor.Start()
}

func (this *Player) Update() {

}

func (this *Player) InitMessage() {
	//玩家登录
	this.RegisterCall("Login", this.HandeLogin)
	//玩家登录到游戏
	this.RegisterCall("C_W_Game_LoginRequset", this.Handle_CW_LoginRequset)
	//创建玩家
	this.RegisterCall("C_W_CreatePlayerRequest", this.Handle_C_W_CreatePlayerRequest)
	//account创建玩家反馈
	this.RegisterCall("CreatePlayer", this.HandleCreatePlayer)
	//玩家断开链接
	this.RegisterCall("Logout", this.HandleLogout)
}

//玩家登录
func (this *Player) HandeLogin(socketId int) {
	PlayerSimpleList := LoadSimplePlayerDatas(this.AccountId)
	this.PlayerSimpleDataList = PlayerSimpleList

	PlayerDataList := make([]*message.PlayerData, len(PlayerSimpleList))
	this.PlayerIdList = []int64{}
	for i, v := range PlayerSimpleList {
		PlayerDataList[i] = &message.PlayerData{PlayerID: proto.Int64(v.PlayerId), PlayerName: proto.String(v.PlayerName), PlayerGold: proto.Int32(int32(v.Gold))}
		this.PlayerIdList = append(this.PlayerIdList, v.PlayerId)
	}

	this.m_Log.Println("玩家登录成功")
	this.SocketId = socketId
	world.SendToClient(socketId, &message.W_C_SelectPlayerResponse{PacketHead: message.BuildPacketHead(this.AccountId, int(message.SERVICE_CLIENT)),
		AccountId:  proto.Int64(this.AccountId),
		PlayerData: PlayerDataList,
	})
}

//玩家登录到游戏
func (this *Player) Handle_CW_LoginRequset(packet *message.C_W_Game_LoginRequset) {
	nPlayerId := packet.GetPlayerId()
	if !this.SetPlayerId(nPlayerId) {
		this.m_Log.Printf("帐号[%d]登入的玩家[%d]不存在", this.AccountId, nPlayerId)
	}

	//读取玩家数据
	this.AddMap()
	//添加到世界频道
	actor.MGR.SendMsg("chatmgr", "AddPlayerToChannel", this.AccountId, this.GetPlayerId(), int64(-3000), this.GetPlayerName(), this.SocketId)
}

//创建玩家
func (this *Player) Handle_C_W_CreatePlayerRequest(packet *message.C_W_CreatePlayerRequest) {
	rows, err := this.m_db.Query(fmt.Sprintf("call `sp_checkcreatePlayer`(%d)", this.AccountId))
	if err == nil && rows != nil {
		if rows.NextResultSet() {
			rs := db.Query(rows, err)
			if rs.Next() {
				err := rs.Row().Int("@err")
				//register
				if err == 0 {
					world.SERVER.GetAccountCluster().BalacaceMsg("W_A_CreatePlayer", this.AccountId, packet.GetPlayerName(), packet.GetSex(), this.GetSocketId())
				} else {
					this.m_Log.Printf("账号[%d]创建玩家上限", this.AccountId)
					world.SendToClient(this.GetSocketId(), &message.W_C_CreatePlayerResponse{
						PacketHead: message.BuildPacketHead(this.AccountId, 0),
						Error:      proto.Int32(int32(err)),
						PlayerId:   proto.Int64(0),
					})
				}
			}
		}
	}
}

//account创建玩家反馈
func (this *Player) HandleCreatePlayer(playerId int64, socketId int, err int) {
	//创建成功
	if err == 0 {
		this.PlayerIdList = []int64{}
		playerSimpleData := LoadSimplePlayerData(playerId)
		this.PlayerSimpleDataList = append(this.PlayerSimpleDataList, playerSimpleData)
		this.PlayerIdList = append(this.PlayerIdList, playerId)
	}

	world.SendToClient(socketId, &message.W_C_CreatePlayerResponse{
		PacketHead: message.BuildPacketHead(this.AccountId, 0),
		Error:      proto.Int32(int32(err)),
		PlayerId:   proto.Int64(playerId),
	})
}

//玩家断开链接
func (this *Player) HandleLogout(accountId int64) {
	world.SERVER.GetLog().Printf("[%d] 断开链接", accountId)
	this.SocketId = 0
	this.Stop()
	this.LeaveMap()
}
