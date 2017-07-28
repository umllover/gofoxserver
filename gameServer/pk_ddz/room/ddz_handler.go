package room

func RegisterHandler(r *DDZ_Entry) {
	r.RegisterBaseFunc()

	r.GetChanRPC().Register("CallScore", r.CallScore)
	r.GetChanRPC().Register("OutCard", r.OutCard)
	r.GetChanRPC().Register("Trustee", r.CTrustee)
	r.GetChanRPC().Register("ShowCard", r.ShowCard)

}
