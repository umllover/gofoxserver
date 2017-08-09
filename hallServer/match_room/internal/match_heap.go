package internal

import (
	"container/heap"
	"errors"
	"fmt"
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

type Item struct {
	value    string // 优先级队列中的数据，可以是任意类型，这里使用string
	priority int    // 优先级队列中节点的优先级
	index    int    // index是该节点在堆中的位置
}

// 优先级队列需要实现heap的interface
type MatchQueue []*Item

// 绑定Len方法
func (pq MatchQueue) Len() int {
	return len(pq)
}

// 绑定Less方法，这里用的是小于号，生成的是小根堆
func (pq MatchQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

// 绑定swap方法
func (pq MatchQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// 绑定put方法，将index置为-1是为了标识该数据已经出了优先级队列了
func (pq *MatchQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	item.index = -1
	return item
}

// 绑定push方法
func (pq *MatchQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// 更新修改了优先级和值的item在优先级队列中的位置
func (pq *MatchQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func main() {
	// 创建节点并设计他们的优先级
	items := map[string]int{"二毛": 5, "张三": 3, "狗蛋": 9}
	i := 0
	pq := make(MatchQueue, len(items)) // 创建优先级队列，并初始化
	for k, v := range items {          // 将节点放到优先级队列中
		pq[i] = &Item{
			value:    k,
			priority: v,
			index:    i}
		i++
	}
	heap.Init(&pq) // 初始化堆
	item := &Item{ // 创建一个item
		value:    "李四",
		priority: 1,
	}
	heap.Push(&pq, item)           // 入优先级队列
	pq.update(item, item.value, 6) // 更新item的优先级
	for len(pq) > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s index:%.2d\n", item.priority, item.value, item.index)
	}
}

//非协程安全 todo 后期改为redis
func IncRoomCnt(roomid int) (int, error) {
	//ret := db.RdsDB.Incr(fmt.Sprintf("IncRoom:%d", roomid))
	//cnt, err := ret.Result()
	//return int(cnt), err
	_, err := db.DB.Exec("UPDATE create_room_info set user_cnt = user_cnt+1 WHERE room_id=?", roomid)
	if err != nil {
		return 0, err
	}

	var Ret []int
	db.DB.Select(&Ret, "SELECT user_cnt FROM create_room_info WHERE room_id=?;", roomid)

	if err != nil {
		return 0, err
	}

	if len(Ret) < 1 {
		return 0, errors.New("not foud error")
	}

	return Ret[0], nil
}

func UpRoomCnt(roomid int, cnt int) {
	_, err := db.DB.Exec("UPDATE create_room_info set user_cnt = ? WHERE room_id=?", cnt, roomid)
	if err != nil {
		log.Error("UpRoomCnt error:%s", err.Error())
	}
}
