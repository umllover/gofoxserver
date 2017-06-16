package common

import (
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"strconv"
)

////////////////////////////////////////////
//全局变量
// TODO 增加 默认(错误)值 参数
func getGlobalVar(key string) string {
	if globalVar, ok := base.GlobalVarCache.Get(key); ok {
		return globalVar.V
	}
	return ""
}

func getGlobalVarFloat64(key string) float64 {
	if value := getGlobalVar(key); value != "" {
		v, _ := strconv.ParseFloat(value, 10)
		return v
	}
	return 0
}

func getGlobalVarInt64(key string, val int64) int64 {
	if value := getGlobalVar(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return val
}

func getGlobalVarInt(key string) int {
	if value := getGlobalVar(key); value != "" {
		v, _ := strconv.Atoi(value)
		return v
	}
	return 0
}

func getGlobalVarIntList(key string) []int {
	if value := getGlobalVar(key); value != "" {
		return utils.GetStrIntList(value)
	}
	return nil
}
