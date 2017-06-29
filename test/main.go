package main

import (
	//"mj/test/Room"
	//"time"
	"fmt"
	"regexp"
)

func main() {
	//room := Room.NewRoom()
	//room.RegistFuncs()
	//room.Run()
	//room.Go("Test", "我是参数哦哦哦哦哦哦")
	//room.Go("SitDown", "我是SitDown")
	////////////////////////
	//time.AfterFunc(1000*time.Hour, func() {})
	//Room.Wg.Wait()
	re, _ := regexp.Compile(`(?:)^.* `)
	src := re.ReplaceAllString("192.168.199.141:8080", "0.0.0.0")
	fmt.Println(src)
}
