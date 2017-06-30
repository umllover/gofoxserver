package room

import (
	"mj/gameServer/common/pk_base/NNBaseLogic"
)

func RegisterHandler(r *NNBaseLogic.NN_PK_base) {
	r.GetChanRPC().Register("Sitdown", r.Sitdown)
	r.GetChanRPC().Register("UserStandup", r.UserStandup)
	r.GetChanRPC().Register("GetUserChairInfo", r.GetUserChairInfo)
	r.GetChanRPC().Register("DissumeRoom", r.DissumeRoom)
	r.GetChanRPC().Register("UserReady", r.UserReady)
	r.GetChanRPC().Register("userRelogin", r.UserReLogin)
	r.GetChanRPC().Register("userOffline", r.UserOffline)
	r.GetChanRPC().Register("SetGameOption", r.SetGameOption)


	r.GetChanRPC().Register("CallScore", r.CallScore)
	r.GetChanRPC().Register("AddScore", r.AddScore)
	/*r.GetChanRPC().Register("CallBanker", r.CallBanker)
	r.GetChanRPC().Register("OxCard", r.OxCard)
	r.GetChanRPC().Register("Qiang", r.Qiang)*/


}
