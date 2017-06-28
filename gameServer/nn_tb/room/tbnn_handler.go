package room

import (
	"mj/gameServer/common/pk_base"
)

func RegisterHandler(r *pk_base.PK_base) {
	r.GetChanRPC().Register("Sitdown", r.Sitdown)
	r.GetChanRPC().Register("UserStandup", r.UserStandup)
	r.GetChanRPC().Register("GetUserChairInfo", r.GetUserChairInfo)
	r.GetChanRPC().Register("DissumeRoom", r.DissumeRoom)
	r.GetChanRPC().Register("UserReady", r.UserReady)
	r.GetChanRPC().Register("userRelogin", r.UserReLogin)
	r.GetChanRPC().Register("userOffline", r.UserOffline)
	r.GetChanRPC().Register("SetGameOption", r.SetGameOption)


	r.GetChanRPC().Register("CallScore", r.SetGameOption)
	r.GetChanRPC().Register("AddScore", r.SetGameOption)
	r.GetChanRPC().Register("CallBanker", r.SetGameOption)
	r.GetChanRPC().Register("OxCard", r.SetGameOption)
	r.GetChanRPC().Register("Qiang", r.SetGameOption)

}
