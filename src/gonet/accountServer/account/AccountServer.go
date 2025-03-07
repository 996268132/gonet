package account

import (
	"database/sql"
	"gonet/base"
	"gonet/common/cluster"
	"gonet/db"
	"gonet/message"
	"gonet/network"
	"log"

	"github.com/golang/protobuf/proto"
)

type (
	ServerMgr struct {
		m_pService   *network.ServerSocket
		m_pServerMgr *ServerSocketManager
		m_pActorDB   *sql.DB
		m_Inited     bool
		m_config     base.Config
		m_Log        base.CLog
		m_AccountMgr *AccountMgr
		m_Cluster    *cluster.Service
		m_SnowFlake  *cluster.Snowflake
	}

	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetDB() *sql.DB
		GetLog() *base.CLog
		GetServer() *network.ServerSocket
		GetServerMgr() *ServerSocketManager
		GetAccountMgr() *AccountMgr
	}

	BitStream base.BitStream
)

var (
	AccountNetIP   string
	AccountNetPort string
	WorkID         int
	DB_Server      string
	DB_Name        string
	DB_UserId      string
	DB_Password    string
	EtcdEndpoints  []string
	SERVER         ServerMgr
)

func (this *ServerMgr) Init() bool {
	if this.m_Inited {
		return true
	}

	//初始化log文件
	this.m_Log.Init("account")
	//初始ini配置文件
	this.m_config.Read("SXZ_SERVER.CFG")
	EtcdEndpoints = this.m_config.Get5("Etcd_Cluster", ",")
	AccountNetIP, AccountNetPort = this.m_config.Get2("Account_LANAddress", ":")
	DB_Server = this.m_config.Get("AccountDB_LANIP")
	DB_Name = this.m_config.Get("AccountDB_Name")
	DB_UserId = this.m_config.Get("AccountDB_UserId")
	DB_Password = this.m_config.Get("AccountDB_Password")

	ShowMessage := func() {
		this.m_Log.Println("**********************************************************")
		this.m_Log.Printf("\tAccountServer Version:\t%s", base.BUILD_NO)
		this.m_Log.Printf("\tAccountServerIP(LAN):\t%s:%s", AccountNetIP, AccountNetPort)
		this.m_Log.Printf("\tActorDBServer(LAN):\t%s", DB_Server)
		this.m_Log.Printf("\tActorDBName:\t\t%s", DB_Name)
		this.m_Log.Println("**********************************************************")
	}
	ShowMessage()

	this.m_Log.Println("正在初始化数据库连接...")
	if this.InitDB() {
		this.m_Log.Printf("[%s]数据库连接是失败...", DB_Name)
		log.Fatalf("[%s]数据库连接是失败...", DB_Name)
		return false
	}
	this.m_Log.Printf("[%s]数据库初始化成功!", DB_Name)

	//账号管理类
	this.m_AccountMgr = new(AccountMgr)
	this.m_AccountMgr.Init(1000)

	//socket管理
	this.m_pServerMgr = new(ServerSocketManager)
	this.m_pServerMgr.Init(1000)

	//初始化socket
	this.m_pService = new(network.ServerSocket)
	port := base.Int(AccountNetPort)
	this.m_pService.Init(AccountNetIP, port)
	this.m_pService.SetMaxReceiveBufferSize(1024)
	this.m_pService.SetMaxSendBufferSize(1024)
	this.m_pService.Start()
	var packet EventProcess
	packet.Init(1000)
	this.m_pService.BindPacketFunc(packet.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_AccountMgr.PacketFunc)
	this.m_pService.BindPacketFunc(this.m_pServerMgr.PacketFunc)

	//snowflake
	//this.m_SnowFlake = cluster.NewSnowflake(UserNetIP, base.Int(UserNetPort), EtcdEndpoints)

	//注册account集群
	//this.m_Cluster = cluster.NewService(int(message.SERVICE_ACCOUNTSERVER), UserNetIP, base.Int(UserNetPort), EtcdEndpoints)
	return false
}

func (this *ServerMgr) InitDB() bool {
	this.m_pActorDB = db.OpenDB(DB_Server, DB_UserId, DB_Password, DB_Name)
	err := this.m_pActorDB.Ping()
	return err != nil
}

func (this *ServerMgr) GetDB() *sql.DB {
	return this.m_pActorDB
}

func (this *ServerMgr) GetLog() *base.CLog {
	return &this.m_Log
}

func (this *ServerMgr) GetServer() *network.ServerSocket {
	return this.m_pService
}

func (this *ServerMgr) GetServerMgr() *ServerSocketManager {
	return this.m_pServerMgr
}

func (this *ServerMgr) GetAccountMgr() *AccountMgr {
	return this.m_AccountMgr
}

func SendToClient(socketId int, packet proto.Message) {
	bitstream := base.NewBitStream(make([]byte, 1024), 1024)
	if !message.GetMessagePacket(packet, bitstream) {
		return
	}
	SERVER.GetServer().SendByID(socketId, bitstream.GetBuffer())
}
