package room

import (
	"math"
	"mj/common/msg"
	. "mj/gameServer/common/mj_logic_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

type GameLogic struct {
	*BaseLogic
	//CardDataArray []int //扑克数据
}

func NewGameLogic() *GameLogic {
	g := new(GameLogic)
	return g
}





