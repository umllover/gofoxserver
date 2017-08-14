package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	. "mj/common/cost"
	"mj/hallServer/common"
	"mj/hallServer/http_service"
	"mj/hallServer/user"
	"net/http"
	"net/url"

	"github.com/lovelly/leaf/log"
)

func (m *UserModule) GetUser(args []interface{}) (interface{}, error) {
	u, ok := m.a.UserData().(*user.User)
	if !ok {
		return nil, errors.New("not foud user Data at GetUser")
	}
	return u, nil
}
func ReqGetMaskCode(phome string, maskCode int) {
	http_service.PostJSON("https://sms.yunpian.com/v2/sms/single_send.json", map[string]interface{}{
		"apikey": common.GetGlobalVar("YUN_PIAN_API_KEY"),
		"mobile": phome,
		"text":   fmt.Sprintf(common.GetGlobalVar(MASK_CODE_TEXT), maskCode),
	})
}

// bingone
func VerifyCode(number string, codes int) int {

	// 修改为您的apikey(https://www.yunpian.com)登录官网后获取
	apikey := common.GetGlobalVar("YUN_PIAN_API_KEY")
	// 修改为您要发送的手机号码，多个号码用逗号隔开
	mobile := number
	// 发送内容
	//text := fmt.Sprintf("【噜噜棋牌游戏中心】您的验证码是%d", codes)
	text := fmt.Sprintf(common.GetGlobalVar(MASK_CODE_TEXT), codes)

	url_send_sms := "https://sms.yunpian.com/v2/sms/single_send.json"

	data_send_sms := url.Values{"apikey": {apikey}, "mobile": {mobile}, "text": {text}}
	log.Debug("data:", data_send_sms)

	return httpsPostForm(url_send_sms, data_send_sms)
}

func httpsPostForm(url string, data url.Values) int {
	resp, err := http.PostForm(url, data)
	defer resp.Body.Close()
	if err != nil {
		log.Debug("error :%s", err.Error())
		return 1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug("error read :%s", err.Error())
		return 2
	}

	fmt.Println(string(body))
	return 0
}
