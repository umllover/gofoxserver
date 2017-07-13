package room

func RegisterHandler(r *ZP_base) {
	r.GetChanRPC().Register("Sitdown", r.Sitdown)
	r.GetChanRPC().Register("UserStandup", r.UserStandup)
	r.GetChanRPC().Register("GetUserChairInfo", r.GetUserChairInfo)
	r.GetChanRPC().Register("DissumeRoom", r.DissumeRoom)
	r.GetChanRPC().Register("UserReady", r.UserReady)
	r.GetChanRPC().Register("userRelogin", r.UserReLogin)
	r.GetChanRPC().Register("userOffline", r.UserOffline)
	r.GetChanRPC().Register("SetGameOption", r.SetGameOption)

	r.GetChanRPC().Register("OutCard", r.OutCard)
	r.GetChanRPC().Register("OperateCard", r.UserOperateCard)
	r.GetChanRPC().Register("SetChaHua", r.ChaHuaMsg)
	r.GetChanRPC().Register("SetBuHua", r.OnUserReplaceCardMsg)
	r.GetChanRPC().Register("SetTingCard", r.OnUserListenCardMsg)
	r.GetChanRPC().Register("UserTrustee", r.OnRecUserTrustee)
}
