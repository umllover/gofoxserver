package user

import (
	"fmt"
	"math"
	"mj/common/msg"
	"mj/hallServer/common"
	"mj/hallServer/db"
	"mj/hallServer/db/model"
	"time"

	"mj/hallServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

//表名字
const (
	day_time_table  = "user_dat_times"
	week_time_table = "week_times"
	time_table      = "user_times"
)

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

	now := time.Now()
	if derr == nil {
		for _, v := range daytimes {
			if !(v.CreateTime.Year() == now.Year() && v.CreateTime.Month() == now.Month() && v.CreateTime.Day() != now.Day()) { //创建时间不是今天
				model.UserDayTimesOp.Delete(u.Id, v.KeyId)
				continue
			}
			u.DayTimes[v.KeyId] = v.V
		}
	}

	//每周次数
	weektimes, werr := model.UserDayTimesOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})

	ny, nd := now.ISOWeek()
	if werr == nil {
		for _, v := range weektimes {
			y, d := v.CreateTime.ISOWeek()
			if !(ny == y && nd == d) { //创建时间不是一周
				model.UserWeekTimesOp.Delete(u.Id, v.KeyId)
				continue
			}
			u.WeekTimes[v.KeyId] = v.V
		}
	}
}

//永久次数
func (u *User) GetForeverTimes(k int) int64 {
	u.RLock()
	defer u.RUnlock()
	return u.Times[k]
}

func (u *User) SetForeverTimes(k int, v int64) {
	u.Lock()
	u.Times[k] = v
	u.Unlock()
	updateTimes(time_table, u.Id, k, v)
}

func (u *User) IncreaseTimes(k int, addv int64) {
	u.Lock()
	u.Times[k] += addv
	u.Unlock()
	updateTimes(time_table, u.Id, k, u.Times[k])
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
	u.RLock()
	defer u.RUnlock()
	return u.DayTimes[k]
}

func (u *User) SetDayTimes(k int, v int64) {
	u.Lock()
	u.DayTimes[k] = v
	u.Unlock()
	updateTimes(day_time_table, u.Id, k, v)
}

func (u *User) IncreaseDayTimes(k int, addv int64) {
	u.Lock()
	u.DayTimes[k] += addv
	u.Unlock()
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
	u.RLock()
	defer u.RUnlock()
	return u.WeekTimes[k]
}

func (u *User) SetWeekTimes(k int, v int64) {
	u.Lock()
	u.WeekTimes[k] = v
	u.Unlock()
	updateTimes(week_time_table, u.Id, k, v)
}

func (u *User) IncreaseWeekTimes(k int, addv int64) {
	u.Lock()
	u.WeekTimes[k] += addv
	u.Unlock()
	updateTimes(week_time_table, u.Id, k, u.WeekTimes[k])
}

func (u *User) GetWeekTimesAll() (data map[int]int64) {
	data = make(map[int]int64)
	for k, v := range u.WeekTimes {
		data[k] = v
	}
	return
}

func (u *User) GetTimes(id int) int64 {
	t, ok := base.ActivityCache.Get(id)
	if !ok {
		log.Error("at GetTimes not foud type :%d", id)
		return math.MaxInt64
	}
	switch t.DrawType {
	case common.ActivityTypeForever:
		return u.GetForeverTimes(id)
	case common.ActivityTypeDay:
		return u.GetDayTimes(id)
	case common.ActivityTypeWeek:
		return u.GetWeekTimes(id)
	}
	return math.MaxInt64
}

func (u *User) HasTimes(id int) bool {
	t, ok := base.ActivityCache.Get(id)
	if !ok {
		log.Error("at GetTimes not foud type :%d", id)
		return false
	}
	switch t.DrawType {
	case common.ActivityTypeForever:
		_, ok := u.Times[id]
		return ok
	case common.ActivityTypeDay:
		_, ok := u.DayTimes[id]
		return ok
	case common.ActivityTypeWeek:
		_, ok := u.WeekTimes[id]
		return ok
	}
	return false
}

func (u *User) SetTimes(id int, v int64) {
	t, ok := base.ActivityCache.Get(id)
	if !ok {
		log.Error("at SetTimes not foud type :%d", id)
		return
	}
	switch t.DrawType {
	case common.ActivityTypeForever:
		u.SetForeverTimes(id, v)
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

//////////////////////////////////////

func (u *User) ClearDayTimes() {
	u.Lock()
	u.DayTimes = make(map[int]int64)
	u.Unlock()
	ClearTimes(day_time_table, u.Id)
}

func (u *User) ClearWeekTimes() {
	u.Lock()
	u.WeekTimes = make(map[int]int64)
	u.Unlock()
	ClearTimes(week_time_table, u.Id)
}

//发送活动次数信息
func (u *User) SendActivityInfo() {
	retMsg := &msg.L2C_ActivityInfo{}
	retMsg.DayTimes = u.GetDayTimrsAll()
	retMsg.Times = u.GetTimrsAll()
	retMsg.WeekTimes = u.GetWeekTimesAll()
	u.WriteMsg(retMsg)
}

func updateTimes(table_name string, uid int64, k int, v int64) bool {
	sql := fmt.Sprintf("insert into %s values(%d,%s,%d) on duplicate key update v=%d;", table_name, uid, k, v, v)
	_, err := db.DB.Exec(sql)
	if err != nil {
		log.Error("at updateTimes error:%s", err.Error())
		return false
	}
	return true
}

func ClearTimes(table_name string, id int64) {
	sql := fmt.Sprintf("delete from %s where user_id=%d;", id)
	_, err := db.DB.Exec(sql)
	if err != nil {
		log.Error("at updateTimes error:%s", err.Error())
		return
	}
}

func ClearTimesByKeys(table_name string, Uid, key int) {
	sql := fmt.Sprintf("delete from %s where user_id=%d and key_id=%d;", Uid, key)
	_, err := db.DB.Exec(sql)
	if err != nil {
		log.Error("at updateTimes error:%s", err.Error())
		return
	}
}
