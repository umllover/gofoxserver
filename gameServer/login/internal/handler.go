package internal

import (
	"mj/common/msg"
	"reflect"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"mj/gameServer/db/model"
	"mj/gameServer/user"
	. "mj/common/cost"
)

var (
	Users = make(map[int]gate.Agent) //key is userId
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	//rpc
	handleRpc("NewAgent", NewAgent, chanrpc.FuncCommon)
	handleRpc("CloseAgent", CloseAgent, chanrpc.FuncCommon)

	//c2s
	handlerC2S(&msg.C2G_GR_LogonMobile{}, handleMBLogin)
}

func handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_GR_LogonMobile)
	retMsg := &msg.G2C_LogonSuccess{}
	agent := args[1].(gate.Agent)
	retcode := 0
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(&msg.G2C_LogonFailur{ResultCode: retcode, DescribeString: "登录失败"})
		} else {
			agent.WriteMsg(retMsg)
		}
	}()

	if recvMsg.UserID == 0 {
		retcode = ParamError
		return
	}

	accountData, ok := model.AccountsinfoOp.Get(recvMsg.UserID)
	if !ok || accountData == nil {
		retcode = NotFoudAccout
		return
	}

	//if accountData.PasswordID != recvMsg.Password {
	//	return
	//}

	user := user.NewUser(accountData.UserID)
	user.Agent = agent
	user.Accountsinfo = accountData
	user.Id = accountData.UserID
	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	agent.SetUserData(user)
}

//连接进来的通知
func NewAgent(args []interface{}) error {
	log.Debug("at game NewAgent")
	return nil
}

//连接关闭的同喜
func CloseAgent (args []interface{}) error {
	log.Debug("at game CloseAgent")
	return nil
}


///////
func loadUser(u *user.User) ( bool){
	ainfo, aok := model.AccountsmemberOp.Get(u.Id, u.Accountsinfo.MemberOrder)
	if !aok {
		log.Error("at loadUser not foud AccountsmemberOp by user", u.Id)
		return false
	}

	log.Debug("load user : == %v", ainfo)
	u.Accountsmember = ainfo

	glInfo, glok := model.GamescorelockerOp.Get(u.Id)
	if !glok {
		log.Error("at loadUser not foud GamescorelockerOp by user %d", u.Id)
		return  false
	}
	u.Gamescorelocker = glInfo

	giInfom, giok := model.GamescoreinfoOp.Get(u.Id)
	if !giok {
		log.Error("at loadUser not foud GamescoreinfoOp by user  %d", u.Id)
		return  false
	}
	u.Gamescoreinfo = giInfom

	ucInfo, uok := model.UserattrOp.Get(u.Id)
	if !uok {
		log.Error("at loadUser not foud UserroomcardOp by user  %d", u.Id)
		return  false
	}
	u.Userattr = ucInfo

	uextInfo, ueok := model.UserextrainfoOp.Get(u.Id)
	if !ueok {
		log.Error("at loadUser not foud UserextrainfoOp by user  %d", u.Id)
		return  false
	}
	u.Userextrainfo = uextInfo
	return  true
}
