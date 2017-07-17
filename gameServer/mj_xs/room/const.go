package room

import "github.com/lovelly/leaf/log"

func GetCardWordArray(index []int) bool {
	CardWordArray := []string{
		"一筒", "二筒", "三筒", "四筒", "五筒", "六筒", "七筒", "八筒", "九筒",
		"一万", "二万", "三万", "四万", "五万", "六万", "七万", "八万", "九万",
		"一条", "二条", "三条", "四条", "五条", "六条", "七条", "八条", "九条",
		"东", "南", "西", "北", "中", "发", "白",
		"春", "夏", "秋", "冬", "梅", "兰", "竹", "菊",
	}

	var data string
	//var data2 []int
	for k, v := range index {
		if v > 0 {
			for i := 0; i < v; i++ {
				data = data + "," + CardWordArray[k]
				//data2 = append(data2, k)
			}
		}
	}
	log.Debug("手牌：%s", data)
	//log.Debug("手牌:%d", len(data2))
	return true
}
