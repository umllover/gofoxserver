package user

import (
	"fmt"
	"mj/hallServer/db"

	"mj/gameServer/db/model"

	"github.com/lovelly/leaf/log"
)


//表名字
const (
	day_time_table = "user_dat_times"
	time_table     = "user_times"
)

//非协程安全
func (u *User) LoadTimes() {
	//永久次数
	times, err := model.UserTimesOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})
	if err == nil {
		for _, v := range times {
			u.Times[v.KeyName] = v.V
		}
	}

	//每日次数
	daytimes, derr := model.UserDayTimesOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})
	if derr == nil {
		for _, v := range daytimes {
			u.DayTimes[v.KeyName] = v.V
		}
	}
}

//每日次数
func (u *User) GetDayTimes(k string) int64 {
	return u.DayTimes[k]
}

func (u *User) SetDatTimes(k string, v int64) {
	u.DayTimes[k] = v
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseDayTimes(k string, addv int64) {
	u.DayTimes[k] += addv
	updateTimes(day_time_table, u.Id, k, u.DayTimes[k])
}

//永久次数
func (u *User) GetTimes(k string) int64 {
	return u.Times[k]
}

func (u *User) SetTimes(k string, v int64) {
	u.Times[k] = v
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseTimes(k string, addv int64) {
	u.Times[k] += addv
	updateTimes(day_time_table, u.Id, k, u.Times[k])
}

func updateTimes(table_name string, uid int, k string, v int64) bool {
	sql := fmt.Sprintf("insert into %s values(%d,%s,%d) on duplicate key update v=%d;", table_name, uid, k, v, v)
	_, err := db.DB.Exec(sql)
	if err != nil {
		log.Error("at updateTimes error:%s", err.Error())
		return false
	}
	return true
}

func ClearTimes(table_name string) {
	sql := fmt.Sprintf("delete from %s")
	_, err := db.DB.Exec(sql)
	if err != nil {
		log.Error("at updateTimes error:%s", err.Error())
		return
	}
}
