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
