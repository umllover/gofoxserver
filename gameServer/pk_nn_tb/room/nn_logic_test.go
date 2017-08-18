package room

import "testing"
import "mj/gameServer/common/pk/pk_base"

func TestCompareCard(t *testing.T) {
	firstData := []int {
		13, 40, 55, 51, 50,
	}
	nextData := []int {
		13, 43, 41, 39, 36,
	}
	nn_logic := NewNNTBZLogic(pk_base.IDX_TBNN)
	if nn_logic.CompareCard(firstData, nextData) {
		t.Error("first data > next data is error", firstData, nextData)
	}

}


func TestCompareCard1(t *testing.T) {
	firstData := []int {
		13, 11, 54, 51, 49,
	}
	nextData := []int {
		13, 42, 39, 7, 54,
	}
	nn_logic := NewNNTBZLogic(pk_base.IDX_TBNN)
	if nn_logic.CompareCard(firstData, nextData) {
	} else {
		t.Error("first data < next data is error", firstData, nextData)
	}
}


func TestCompareCard2(t *testing.T) {
	bankerData := []int { 45, 28, 25, 22, 5}
	playerData := []int { 61, 9, 40, 2, 1}
	nn_logic := NewNNTBZLogic(pk_base.IDX_TBNN)
	if nn_logic.CompareCard(bankerData, playerData) {
		t.Error( "banker win is error ", bankerData, playerData)
	} 
}

