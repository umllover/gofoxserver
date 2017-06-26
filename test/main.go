package main

import (
	"mj/test/Room"
	"time"
)

func main() {
	room := Room.NewRoom()
	room.RegistFuncs()
	room.Run()
	room.Go("Test", "我是参数哦哦哦哦哦哦")
	room.Go("SitDown", "我是SitDown")
	//////////////////////
	time.AfterFunc(1000*time.Hour, func() {})
	Room.Wg.Wait()
}
