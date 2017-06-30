package main

import (
	//"mj/test/Room"
	//"time"
	"fmt"
)

func tt() (int, int) {
	return 1, 2
}
func main() {
	//room := Room.NewRoom()
	//room.RegistFuncs()
	//room.Run()
	//room.Go("Test", "我是参数哦哦哦哦哦哦")
	//room.Go("SitDown", "我是SitDown")
	////////////////////////
	//time.AfterFunc(1000*time.Hour, func() {})
	//Room.Wg.Wait()
	tt := 0
	foo := func() {
		tt++
	}

	foo()
	fmt.Println(tt)

}
