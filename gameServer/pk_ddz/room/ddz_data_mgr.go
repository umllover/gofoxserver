package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
	CurrentUser  int              // 当前玩家
	OutCardCount [GAME_PLAYER]int // 出牌次数
	WaitTime     int              // 等待时间
	TimerControl int              // 时间控制

	// 托管信息
	OffLineTrustee bool // 离线托管

	// 炸弹信息
	BombCount     int              // 炸弹个数
	EachBombCount [GAME_PLAYER]int // 炸弹个数

	// 叫分信息
	CallScoreCount int              // 叫分次数
	BankerScore    int              // 庄家叫分
	ScoreInfo      [GAME_PLAYER]int // 叫分信息

	// 出牌信息
	TurnWiner     int            // 胜利玩家
	TurnCardCount int            // 出牌数目
	TurnCardData  [MAX_COUNT]int // 出牌数据

	// 扑克信息
	BankerCard    [3]int                      // 游戏底牌
	HandCardCount [GAME_PLAYER]int            // 扑克数目
	HandCardData  [GAME_PLAYER][MAX_COUNT]int // 手上扑克

	// 组件变量
	//CGameLogic						m_GameLogic;						//游戏逻辑

	// 组件接口
	//ITableFrame	*					m_pITableFrame;						//框架接口
	//tagCustomRule *					m_pGameCustomRule;					//自定规则
	//tagGameServiceOption *			m_pGameServiceOption;				//游戏配置
	//tagGameServiceAttrib *			m_pGameServiceAttrib;				//游戏属性
}

// 游戏开始
func (room *DDZ_Entry) OnEventGameStart() {

}
