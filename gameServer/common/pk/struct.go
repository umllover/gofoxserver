package pk

//分析结构
type TagAnalyseType struct {
	BOnePare        bool   //有一对
	BTwoPare        bool   //有两对
	BThreeSame      bool   //有三条
	BStraight       bool   //有顺子
	BFlush          bool   //有同花
	BGourd          bool   //有葫芦
	BFourSame       bool   //有铁支
	BStraightFlush  bool   //有同花顺
	CbOnePare       []int  //一对的序号
	CbTwoPare       []int  //两对的序号
	CbThreeSame     []int  //三条的序号
	CbStraight      []int  //顺子的序号
	CbFlush         []int  //同花的序号
	CbGourd         []int  //葫芦的序号
	CbFourSame      []int  //铁支的序号
	CbStraightFlush []int  //同花顺的序号
	BbOnePare       []bool //所有一对标志弹起
	BbTwoPare       []bool //所有二对标志弹起
	BbThreeSame     []bool //所有三条标志弹起
	BbStraight      []bool //所有顺子标志弹起
	BbFlush         []bool //所有同花标志弹起
	BbGourd         []bool //所有葫芦标志弹起
	BbFourSame      []bool //所有铁支标志弹起
	BbStraightFlush []bool //所有同花顺标志弹起
	BtOnePare       int    //一对的数量 单独
	BtThreeSame     int    //三条数量   单独

	BtTwoPare       int //两对的数量
	BtStraight      int //顺子的数量
	BtFlush         int //同花的数量
	BtGourd         int //葫芦的数量
	BtFourSame      int //铁支的数量
	BtStraightFlush int //同花顺的数量
}
