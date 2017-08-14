package room

func RegisterHandler(r *SSS_Entry) {
	r.RegisterBaseFunc()
	r.GetChanRPC().Register("ShowCard", r.ShowSSSCard)
	r.GetChanRPC().Register("Trustee", r.Trustee)
}
