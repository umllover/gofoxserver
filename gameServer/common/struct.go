package common

//游戏房间
type tagGameServer struct {
	wKindID int						//名称索引
	wNodeID int							//节点索引
	wSortID int							//排序索引
	wServerID int							//房间索引
	wServerPort int						//房间端口
	wServerType int						//房间类型
	DdwOnLineCount int					//在线人数
	DdwFullCount int						//满员人数

	//fdl add 手机房间显示上下限
	lRestrictScore int						//限制积分
	lMinTableScore int						//最低积分
	lMinEnterScore int						//最低积分
	lMaxEnterScore int						//最高积分

	//fdl add end
	szServerAddr string					//房间名称
	szServerName string					//房间名称
	//私人房添加
	dwSurportType int					//支持类型
	wTableCount int						//桌子数目
};
