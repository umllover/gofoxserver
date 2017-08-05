package conf

import (
	"log"
	"time"
)

var (
	// log conf
	LogFlag = log.LstdFlags | log.Llongfile

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 65535
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 2
	LittleEndian           = false

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000

	// agent conf
	AgentGoLen              = 50
	AgentTimerDispatcherLen = 50
	AgentAsynCallLen        = 50
	AgentChanRPCLen         = 50

	// cluster conf
	HeartBeatInterval = 5

	// room
	MaxRoomMsgLen       = 50
	DestroyRoomInterval = 3600

	Shutdown = false
	Test     = false
)
