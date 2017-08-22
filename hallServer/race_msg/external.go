package race_msg

import (
	"mj/hallServer/race_msg/internal"
)

var (
	ChanRPC  = internal.ChanRPC
	Module   = new(internal.Module)
	GmNotify = internal.ReciveGMMsg
)
