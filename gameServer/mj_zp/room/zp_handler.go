package room

func RegisterHandler(r *ZP_base) {
	r.RegisterBaseFunc()

	r.GetChanRPC().Register("OutCard", r.OutCard)
	r.GetChanRPC().Register("OperateCard", r.UserOperateCard)
	r.GetChanRPC().Register("SetChaHua", r.ChaHuaMsg)
	r.GetChanRPC().Register("SetBuHua", r.OnUserReplaceCardMsg)
	r.GetChanRPC().Register("SetTingCard", r.OnUserListenCardMsg)
	r.GetChanRPC().Register("UserTrustee", r.OnRecUserTrustee)
}
