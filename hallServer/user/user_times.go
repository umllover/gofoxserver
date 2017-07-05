package user

import (
	"fmt"
	"math"
	"mj/hallServer/common"
	"mj/hallServer/db"
	"mj/hallServer/db/model"

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
			u.Times[v.KeyId] = v.V
		}
	}

	//每日次数
	daytimes, derr := model.UserDayTimesOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})

	if derr == nil {
		for _, v := range daytimes {
			u.DayTimes[v.KeyId] = v.V
		}
	}

}

//永久次数
func (u *User) GetTimes(k int) int64 {
	return u.Times[k]
}

func (u *User) SetTimes(k int, v int64) {
	u.Times[k] = v
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseTimes(k int, addv int64) {
	u.Times[k] += addv
	updateTimes(day_time_table, u.Id, k, u.Times[k])
}

func (u *User) GetTimrsAll() (data map[int]int64) {
	data = make(map[int]int64)
	for k, v := range u.Times {
		data[k] = v
	}
	return
}

//每日次数
func (u *User) GetDayTimes(k int) int64 {
	return u.DayTimes[k]
}

func (u *User) SetDayTimes(k int, v int64) {
	u.DayTimes[k] = v
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseDayTimes(k int, addv int64) {
	u.DayTimes[k] += addv
	updateTimes(day_time_table, u.Id, k, u.DayTimes[k])
}

func (u *User) GetDayTimrsAll() (data map[int]int64) {
	data = make(map[int]int64)
	for k, v := range u.DayTimes {
		data[k] = v
	}
	return
}

//每周次数
func (u *User) GetWeekTimes(k int) int64 {
	return u.WeekTimes[k]
}

func (u *User) SetWeekTimes(k int, v int64) {
	u.WeekTimes[k] = v
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseWeekTimes(k int, addv int64) {
	u.WeekTimes[k] += addv
	updateTimes(day_time_table, u.Id, k, u.WeekTimes[k])
}

func (u *User) GetWeekTimesAll() (data map[int]int64) {
	data = make(map[int]int64)
	for k, v := range u.WeekTimes {
		data[k] = v
	}
	return
}

func (u *User) GetTimesByType(id int, types int) int64 {
	switch types {
	case common.ActivityTypeForever:
		return u.GetTimes(id)
	case common.ActivityTypeDay:
		return u.GetDayTimes(id)
	case common.ActivityTypeWeek:
		return u.GetWeekTimes(id)
	}
	return math.MaxInt64
}

func (u *User) SetTimesByType(id int, v int64, types int) {
	switch types {
	case common.ActivityTypeForever:
		u.SetTimes(id, v)
	case common.ActivityTypeDay:
		u.SetDayTimes(id, v)
	case common.ActivityTypeWeek:
		u.SetWeekTimes(id, v)
	}
}

func (u *User) IncreaseTimesByType(id int, v int64, types int) {
	switch types {
	case common.ActivityTypeForever:
		u.IncreaseTimes(id, v)
	case common.ActivityTypeDay:
		u.IncreaseDayTimes(id, v)
	case common.ActivityTypeWeek:
		u.IncreaseWeekTimes(id, v)
	}
}

func updateTimes(table_name string, uid int, k int, v int64) bool {
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
