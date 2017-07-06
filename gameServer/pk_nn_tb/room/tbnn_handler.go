package room

func RegisterHandler(r *NNTB_Entry) {
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
	r.GetChanRPC().Register("OpenCard", r.OpenCard)
	/*r.GetChanRPC().Register("Banker", r.Banker)
	r.GetChanRPC().Register("OxCard", r.OxCard)
	r.GetChanRPC().Register("Qiang", r.Qiang)*/

}
