package room

func RegisterHandler(r *hz_entry) {
	r.RegisterBaseFunc()

	r.GetChanRPC().Register("OutCard", r.OutCard)
	r.GetChanRPC().Register("OperateCard", r.UserOperateCard)
}
