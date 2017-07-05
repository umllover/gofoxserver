package room

func RegisterHandler(r *hz_entry) {
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
	r.GetChanRPC().Register("C2G_HZMJ_ZhaMa", r.ZhaMa)
}