package room

import (
	"mj/common/msg"
)
type GameLogic struct {

}
var DefaultGameLogic = NewGameLogic()

func NewGameLogic()*GameLogic {
	return new(GameLogic)
}

func (g *GameLogic)AnalyseTingCard(CardIndex []int8,  WeaveItem []*msg.WeaveItem, WeaveCount, OutCardCount int8,OutCardData,HuCardCount []int8, HuCardData[][]int8)(int8){
	return 0
}

func (g *GameLogic)GetCardCount(cCardIndex []int8) int8 {
	return 0
}

func (g *GameLogic)SwitchToCardData(CardIndex, CardData []int8)(int8){
	return 0
}
