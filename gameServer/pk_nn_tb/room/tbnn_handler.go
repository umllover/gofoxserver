package room

func RegisterHandler(r *NNTB_Entry) {
	r.RegisterBaseFunc()

	r.GetChanRPC().Register("CallScore", r.CallScore)
	r.GetChanRPC().Register("AddScore", r.AddScore)
	r.GetChanRPC().Register("OpenCard", r.OpenCard)
	/*r.GetChanRPC().Register("Banker", r.Banker)
	r.GetChanRPC().Register("OxCard", r.OxCard)
	r.GetChanRPC().Register("Qiang", r.Qiang)*/

}
