package gate

import (
	"gonet/actor"
	"gonet/base"
	"sync"
	"time"
)

type (
	PlayerManager struct {
		actor.Actor
		m_SocketMap  map[int]int64
		m_AccountMap map[int64]*AccountInfo
		m_Locker     *sync.RWMutex
	}

	IPlayerMangaer interface {
		actor.IActor
		ReleaseSocketMap(int, bool)
		AddAccountMap(int, int64) int
		GetSocket(int64) int
		GetAccount(int) int64
		GetAccountInfo(int) *AccountInfo
	}

	AccountInfo struct {
		AccountId int64
		LastTime  int64
		SocketId  int
		WSocketId uint32
	}
)

var (
	g_pAccount = &AccountInfo{}
)

func NewAccountInfo(socket int, accountId int64) *AccountInfo {
	accountInfo := AccountInfo{LastTime: time.Now().Unix(), SocketId: socket, WSocketId: 0, AccountId: accountId}
	return &accountInfo
}

func (this *PlayerManager) ReleaseSocketMap(socketId int, bClose bool) {
	this.m_Locker.RLock()
	accountid, _ := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	this.m_Locker.Lock()
	delete(this.m_AccountMap, accountid)
	delete(this.m_SocketMap, socketId)
	this.m_Locker.Unlock()
	//if bClose{
	SERVER.GetServer().StopClient(socketId)
	//}
}

func (this *PlayerManager) AddAccountMap(accountId int64, socketId int) int {
	Id := this.GetSocket(accountId)
	this.ReleaseSocketMap(Id, Id != socketId)

	accountInfo := NewAccountInfo(socketId, accountId)
	accountInfo.WSocketId = SERVER.GetWorldCluster().RandomCluster()
	this.m_Locker.Lock()
	this.m_AccountMap[accountId] = accountInfo
	this.m_SocketMap[socketId] = accountId
	this.m_Locker.Unlock()
	SERVER.GetWorldCluster().SendMsg(accountInfo.WSocketId, "G_W_CLoginRequest", accountId)
	return base.NONE_ERROR
}

func (this *PlayerManager) GetSocket(accountId int64) int {
	socketId := 0
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist {
		socketId = accountInfo.SocketId
	}
	return socketId
}

func (this *PlayerManager) GetAccount(socketId int) int64 {
	accoundId := int64(0)
	this.m_Locker.RLock()
	id, exist := this.m_SocketMap[socketId]
	this.m_Locker.RUnlock()
	if exist {
		accoundId = id
	}
	return accoundId
}

func (this *PlayerManager) GetAccountInfo(socketId int) *AccountInfo {
	accountId := this.GetAccount(socketId)
	this.m_Locker.RLock()
	accountInfo, exist := this.m_AccountMap[accountId]
	this.m_Locker.RUnlock()
	if exist {
		return accountInfo
	}
	return nil
}

func (this *PlayerManager) Init(num int) {
	this.Actor.Init(num)
	this.m_SocketMap = make(map[int]int64)
	this.m_AccountMap = make(map[int64]*AccountInfo)
	this.m_Locker = &sync.RWMutex{}

	this.InitMessage()
	this.Actor.Start()
}

func (this *PlayerManager) InitMessage() {
	this.RegisterCall("ADD_ACCOUNT", this.Handle_ADD_ACCOUNT)

	this.RegisterCall("DEL_ACCOUNT", this.Handle_DEL_ACCOUNT)

	//重连世界服务器，账号重新登录
	this.RegisterCall("Account_Relink", this.Handle_Account_Relink)
}

func (this *PlayerManager) Handle_ADD_ACCOUNT(accountId int64, socketId int) {
	SERVER.GetLog().Printf("login incoming  Socket:%d Account:%d ", socketId, accountId)
	this.AddAccountMap(accountId, socketId)
}

func (this *PlayerManager) Handle_DEL_ACCOUNT(socketid int) {
	accountId := this.GetAccount(socketid)
	this.ReleaseSocketMap(socketid, true)
	SERVER.GetAccountCluster().BoardCastMsg("G_ClientLost", accountId)
}

//重连世界服务器，账号重新登录
func (this *PlayerManager) Handle_Account_Relink() {
	accountMap := make(map[int64]uint32)
	this.m_Locker.RLock()
	for i, v := range this.m_AccountMap {
		accountMap[i] = v.WSocketId
	}
	this.m_Locker.RUnlock()

	if len(accountMap) != 0 {
		for i, v := range accountMap {
			SERVER.GetWorldCluster().SendMsg(v, "G_W_CLoginRequest", i)
		}
	}
}
