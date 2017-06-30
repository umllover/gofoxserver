package proto



// import (
// 	"changit.cn/contra/server/clientSimulator/clientSimulator"
// 	"changit.cn/contra/server/services"
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	//simplejson "github.com/bitly/go-simplejson"
// 	"changit.cn/contra/server/model"
// 	"log"
// 	"testing"
// )

// var host = flag.String("host", "127.0.0.1", "游戏地址")
// var port = flag.Int("port", 9001, "游戏端口")
// var authHost = flag.String("auth_host", "192.168.199.156", "认证服务器地址")
// var authPort = flag.Int("auth_port", 80, "认证服务器端口")

// var mockClient = &clientSimulator.MockClient{}

// func connect() {
// 	serverUrl := fmt.Sprintf("%s:%d", *host, *port)
// 	wsUrl := fmt.Sprintf("ws://%s:%d/rpc", *host, *port)
// 	authUrl := fmt.Sprintf("http://%s:%d/auth", *authHost, *authPort)

// 	mockClient.CheckUser("changtest99", "123456", authUrl)
// 	mockClient.ConnectWS(serverUrl, wsUrl)
// 	logger.Info("connected")
// }

// func TestLogin(t *testing.T) {

// 	connect()

// 	req := &clientSimulator.Request{FuncName: "Login", Params: []interface{}{mockClient.UserId, mockClient.SessionKey, 31}}

// 	_, _ = json.Marshal(req)
// 	//logger.Info("send:", string(reqJson))

// 	mockClient.SendReq(req)

// 	result := mockClient.GetResp()

// 	resp := clientSimulator.Response{}
// 	json.Unmarshal(result, &resp)

// 	//logger.Info("recv:", string(result))
// 	//logger.Info("recv:", resp.Error.Code)

// 	if "Login" == resp.FuncName && 0 == resp.Error.Code {

// 	} else {
// 		t.Error("Login OK")
// 	}
// }

// func TestPlayerInfo(t *testing.T) {

// 	req := &clientSimulator.Request{FuncName: "PlayerInfo", Params: []interface{}{}}

// 	mockClient.SendReq(req)

// 	result := mockClient.GetResp()

// 	respData := map[string]interface{}{}

// 	json.Unmarshal(result, &respData)

// 	fire_set, _ := json.Marshal(respData["data"].(map[string]interface{})["fire_set"])

// 	mockClient.FireSet = &model.PlayerFireSet{}
// 	json.Unmarshal(fire_set, mockClient.FireSet)
// 	//logger.Info("recv:", mockClient.FireSet)

// }

// func TestResetHerTalent(t *testing.T) {

// 	req := &clientSimulator.Request{FuncName: "ResetHerTalent", Params: []interface{}{mockClient.FireSet.HeroId, 0}}

// 	mockClient.SendReq(req)
// 	respData := &clientSimulator.Response{}
// 	json.Unmarshal(mockClient.GetResp(), &respData)

// 	if "" == respData.FuncName && services.ErrTalentUnadd == respData.Error.Code {
// 		logger.Info("天赋未加点")
// 	}

// 	//{"func_name":"PlayerAttr","data":{"stone":2646},"error":{}}
// }
