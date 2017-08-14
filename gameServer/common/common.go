package common

import (
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"strconv"
)

const (
	KIND_TYPE_HZMJ = 389 //红中麻将
	KIND_TYPE_ZPMJ = 391 //漳浦麻将
	KIND_TYPE_TBNN = 28  //同比牛牛
	KIND_TYPE_DDZ  = 29  //斗地主
	KIND_TYPE_SSS  = 30  //十三水
)

const (
	TableFullCount = 1
)

////////////////////////////////////////////
//全局变量
// TODO 增加 默认(错误)值 参数
func GetGlobalVar(key string) string {
	if globalVar, ok := base.GlobalVarCache.Get(key); ok {
		return globalVar.V
	}
	return ""
}

func GetGlobalVarFloat64(key string) float64 {
	if value := GetGlobalVar(key); value != "" {
		v, _ := strconv.ParseFloat(value, 10)
		return v
	}
	return 0
}

func GetGlobalVarInt64(key string, val int64) int64 {
	if value := GetGlobalVar(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return val
}

func GetGlobalVarInt(key string) int {
	if value := GetGlobalVar(key); value != "" {
		v, _ := strconv.Atoi(value)
		return v
	}
	return 0
}

func GetGlobalVarIntList(key string) []int {
	if value := GetGlobalVar(key); value != "" {
		return utils.GetStrIntList(value)
	}
	return nil
}
