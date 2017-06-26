package room

import (
	"math"
	"mj/common/msg"
	. "mj/gameServer/common/mj_logic_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

//花色
const (
	LOGIC_MASK_COLOR = 0xF0
	LOGIC_MASK_VALUE = 0x0F
)

// 扑克类型
const (
	OX_VALUE0 		=			0									//混合牌型
	OX_THREE_SAME	=			105									//小牛牛——5张牌都小于5（含5），并且5张牌相加不大于10
	OX_FOUR_SAME	=			104									////炸弹——5张牌中有4张一样的牌。
	OX_FOURKING		=			102									//天王牌型四花
	OX_FIVEKING		=			103									//天王牌型五花
)


type GameLogic struct {
	*BaseLogic
	//CardDataArray []int //扑克数据
}


//扑克数据
var  CardDataArray  = []int {
	0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,0x09,0x0A,0x0B,0x0C,0x0D,	//方块 A - K
	0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,	//梅花 A - K
	0x21,0x22,0x23,0x24,0x25,0x26,0x27,0x28,0x29,0x2A,0x2B,0x2C,0x2D,	//红桃 A - K
	0x31,0x32,0x33,0x34,0x35,0x36,0x37,0x38,0x39,0x3A,0x3B,0x3C,0x3D,	//黑桃 A - K 3 13 54
	0x4E,0x4F,//14 15
}


func NewGameLogic() *GameLogic {
	g := new(GameLogic)
	return g
}



//分析结构
type   tagAnalyseResult		struct {
	FourCount					int						//四张数目
	ThreeCount					int						//三张数目
	DoubleCount					int						//两张数目
	SignedCount					int						//单张数目
	FourLogicVolue				[]int					//四张列表
	ThreeLogicVolue				[]int					//三张列表
	DoubleLogicVolue			[]int					//两张列表
	SignedLogicVolue			[]int					//单张列表
	FourCardData				[]int					//四张列表
	ThreeCardData				[]int					//三张列表
	DoubleCardData				[]int					//两张列表
	SignedCardData				[]int					//单张数目
}




//获取类型
func (g *GameLogic) GetCardType(cardData []int, cardCount int) (cardType int) {
	cardType = 0
	defer func() {
		if cardType == 0 {
			log.Debug("get card type failed")
		}
	}()

	if cardCount != MAX_COUNT {
		return
	}
	SameCount := 0
	//SortCardList(cbCardData,cbCardCount)

	var Temp [MAX_COUNT]int
	Sum :=0

	for i:=0; i<cardCount; i++{
		Temp[i] = GetCardLogicValue(cardData[i])
		//printf("%d\n",bTemp[i])
		log.Debug(Temp[i])
		Sum += Temp[i]
	}
	//printf("%d\n",bSum)
	log.Debug(Sum)
	//葫芦牌型 改成 小牛牛——5张牌都小于5（含5），并且5张牌相加不大于10。
	//if(bXiaoNiu&&bSum<=10)
	//{
	//	return OX_THREE_SAME //小牛牛——5张牌都小于5（含5），并且5张牌相加不大于10。
	//}


	//BYTE bSecondValue = GetCardValue(cbCardData[MAX_COUNT/2])
	//for(BYTE i=0i<cbCardCounti++)
	//{
	//	if(bSecondValue == GetCardValue(cbCardData[i]))
	//	{
	//		bSameCount++
	//	}
	//}
	//if(bSameCount==4)return OX_FOUR_SAME//炸弹——5张牌中有4张一样的牌。
	//没牛-牛丁-牛二……牛八-牛九-牛牛-四花-五花-炸弹-小牛牛

	//王的数量
	BYTE bKingCount=0,bTenCount=0
	for(BYTE i=0i<cbCardCounti++)
	{
	if(GetCardValue(cbCardData[i])>10 && cbCardData[i]!=0x4E && cbCardData[i]!=0x4F)
	{
	bKingCount++
	}
	//if(GetCardValue(cbCardData[i])>10)
	//if (cbCardData[i]==0x4E||cbCardData[i]==0x4F)
	//{
	//	bKingCount++
	//}
	else if(GetCardValue(cbCardData[i])==10)
	{
	bTenCount++
	}
	}
	if(bKingCount==MAX_COUNT) return OX_FIVEKING//五花——5张牌都是10以上（不含10）的牌。。
	//	else if(bKingCount==MAX_COUNT-1 && bTenCount==1) return OX_FOURKING//四花——5张牌中有1张是10，其余4张是10以上（不含10）的牌

	////葫芦牌型
	//if(bSameCount==3)
	//{
	//	if((bSecondValue!=GetCardValue(cbCardData[3]) && GetCardValue(cbCardData[3])==GetCardValue(cbCardData[4]))
	//	||(bSecondValue!=GetCardValue(cbCardData[1]) && GetCardValue(cbCardData[1])==GetCardValue(cbCardData[0])))
	//		return OX_THREE_SAME
	//}


	BYTE cbValue=GetCardLogicValue(cbCardData[3])
	cbValue+=GetCardLogicValue(cbCardData[4])

	if(cbValue>10)
	{														//防止原值变小
		if( (cbCardData[3]==0x4E||cbCardData[4]==0x4F||cbCardData[4]==0x4E||cbCardData[3]==0x4F) /*&& (cbValue<20)*/)
		{
			cbValue=10
		}
else
cbValue-=10 //2.3
}

//for (BYTE i=0i<cbCardCount-1i++)
//{
//	for (BYTE j=i+1j<cbCardCountj++)
//	{
//		if((bSum-bTemp[i]-bTemp[j])%10==0)
//		{
//			return ((bTemp[i]+bTemp[j])>10)?(bTemp[i]+bTemp[j]-10):(bTemp[i]+bTemp[j])
//		}
//	}
//}

	return  cardType
}

//获取倍数
BYTE CGameLogic::GetTimes(BYTE cbCardData[], BYTE cbCardCount,BYTE bniu)
{
if (bniu!=1)
{
return 1
}
if(cbCardCount!=MAX_COUNT)return 0

BYTE bTimes=GetCardType(cbCardData,MAX_COUNT)
printf("========%d\n",bTimes)
if(bTimes<7)return 1
else if(bTimes==7)return 1
else if(bTimes==8)return 2
else if(bTimes==9)return 3
else if(bTimes==10)return 4
//else if(bTimes==OX_THREE_SAME)return 5
//else if(bTimes==OX_FOUR_SAME)return 5
//else if(bTimes==OX_FOURKING)return 5
else if(bTimes==OX_FIVEKING)return 5
//////BYTE bTimes=0/*GetCardType(cbCardData,MAX_COUNT)*/

//////BYTE cbValue=GetCardLogicValue(cbCardData[3])
//////cbValue+=GetCardLogicValue(cbCardData[4])

/////bool firstking=false
/////bool nextking=false



////if(cbValue>10)
/////{														//防止原值变小
/////	if( (cbCardData[3]==0x4E||cbCardData[4]==0x4F||cbCardData[4]==0x4E||cbCardData[3]==0x4F) /*&& (cbValue<20)*/)
/////	{
/////		nextking=true
/////		cbValue=10
/////	}
/////	else
/////		cbValue-=10 //2.3
/////}
/////bTimes=cbValue
//if(bTimes<8)return 5
////else if(bTimes==7)return 2
//else if(bTimes==8)return 10
//else if(bTimes==9)return 10
//else if(bTimes==10)return 15
//else if(bTimes==OX_THREE_SAME)return 40//小牛牛
//else if(bTimes==OX_FOUR_SAME)return 30//炸弹
//else if(bTimes==OX_FOURKING)return 20
//else if(bTimes==OX_FIVEKING)return 25//
/////	if(bTimes<10)return 1
/////else if(bTimes>=10)return 2
return 0
}

//获取牛牛
//bool CGameLogic::GetOxCard(BYTE cbCardData[], BYTE cbCardCount)
//{
//	ASSERT(cbCardCount==MAX_COUNT)
//
//	//设置变量
//	BYTE bTemp[MAX_COUNT],bTempData[MAX_COUNT]
//	CopyMemory(bTempData,cbCardData,sizeof(bTempData))
//	BYTE bSum=0
//	for (BYTE i=0i<cbCardCounti++)
//	{
//		bTemp[i]=GetCardLogicValue(cbCardData[i])
//		bSum+=bTemp[i]
//	}
//
//	//查找牛牛
//	for (BYTE i=0i<cbCardCount-1i++)
//	{
//		for (BYTE j=i+1j<cbCardCountj++)
//		{
//			if((bSum-bTemp[i]-bTemp[j])%10==0)
//			{
//				BYTE bCount=0
//				for (BYTE k=0k<cbCardCountk++)
//				{
//					if(k!=i && k!=j)
//					{
//						cbCardData[bCount++] = bTempData[k]
//					}
//				}ASSERT(bCount==3)
//
//				cbCardData[bCount++] = bTempData[i]
//				cbCardData[bCount++] = bTempData[j]
//
//				return true
//			}
//		}
//	}
//
//	return false
//}
bool CGameLogic::GetOxCard(BYTE cbCardData[], BYTE cbCardCount)
{
ASSERT(cbCardCount==MAX_COUNT)

//设置变量
BYTE bTemp[MAX_COUNT],bTempData[MAX_COUNT]
CopyMemory(bTempData,cbCardData,sizeof(bTempData))
BYTE bSum=0
for (BYTE i=0i<cbCardCounti++)
{
bTemp[i]=GetCardLogicValue(cbCardData[i])
bSum+=bTemp[i]
}
//王的数量
BYTE bKingCount=0,bTenCount=0
for(BYTE i=0i<cbCardCounti++)
{
/*if(GetCardValue(cbCardData[i])>10)
{
	bKingCount++
}*/
//if(GetCardValue(cbCardData[i])>10)
if (cbCardData[i]==0x4E||cbCardData[i]==0x4F)
{
bKingCount++
}
else if(GetCardValue(cbCardData[i])==10)
{
bTenCount++
}
}
BYTE maxNiuZhi=0
BYTE NiuWeiZhi=0//记录最大牛牌的位置
BYTE bNiuTemp[30][MAX_COUNT]

BYTE bIsKingPai[30]//是否组成牌型时王变了数值
ZeroMemory(bIsKingPai,sizeof(bIsKingPai))

ZeroMemory(bNiuTemp,sizeof(bNiuTemp))

BYTE NiuShu=0
bool bHaveKing=false
//查找牛牛
for (BYTE i=0i<cbCardCount-1i++)
{
for (BYTE j=i+1j<cbCardCountj++)
{
bHaveKing=false
BYTE ShengYu=(bSum-bTemp[i]-bTemp[j])%10
if (ShengYu>0 && bKingCount>0)
{
BYTE bCount=0
for (BYTE k=0k<cbCardCountk++)
{
if(k!=i && k!=j)
{
//bNiuTemp[NiuShu][bCount++] = bTempData[k]
if (bTempData[k]==0x4E||bTempData[k]==0x4F)
{
bHaveKing=true
}
}
}
}

if(( (bSum-bTemp[i]-bTemp[j])%10==0 )||
bHaveKing ) //如果减去2个剩下3个是10的倍数
{
BYTE bCount=0
for (BYTE k=0k<cbCardCountk++)
{
if(k!=i && k!=j)
{
//cbCardData[bCount++] = bTempData[k]
bNiuTemp[NiuShu][bCount++] = bTempData[k]
}
}ASSERT(bCount==3)

bNiuTemp[NiuShu][bCount++] = bTempData[i]
bNiuTemp[NiuShu][bCount++] = bTempData[j]


BYTE cbValue=bTemp[i]
cbValue+=bTemp[j]
if(cbValue>10)
{
if (bTempData[i]==0x4E||bTempData[j]==0x4F||bTempData[i]==0x4F||bTempData[j]==0x4E)
{
bHaveKing=true
cbValue=10
}
else
cbValue-=10 //2.3
}

bIsKingPai[NiuShu]=bHaveKing

if (cbValue>maxNiuZhi)
{
maxNiuZhi = cbValue//最大牛数量
NiuWeiZhi = NiuShu//记录最大牛牌的位置
}
/*cbCardData[bCount++] = bTempData[i]
cbCardData[bCount++] = bTempData[j]*/
NiuShu++
continue
//return true
}
}
}

if (NiuShu>0)
{
for(BYTE i=0i<cbCardCounti++)
{
cbCardData[i]=bNiuTemp[NiuWeiZhi][i]
}

return true
}
return false
}

//获取整数
bool CGameLogic::IsIntValue(BYTE cbCardData[], BYTE cbCardCount)
{
BYTE sum=0
for(BYTE i=0i<cbCardCounti++)
{
sum+=GetCardLogicValue(cbCardData[i])
}
ASSERT(sum>0)
return (sum%10==0)
}

//排列扑克
void CGameLogic::SortCardList(BYTE cbCardData[], BYTE cbCardCount)
{
//转换数值
BYTE cbLogicValue[MAX_COUNT]
for (BYTE i=0i<cbCardCounti++) cbLogicValue[i]=GetCardValue(cbCardData[i])

//排序操作
bool bSorted=true
BYTE cbTempData,bLast=cbCardCount-1
do
{
bSorted=true
for (BYTE i=0i<bLasti++)
{
if ((cbLogicValue[i]<cbLogicValue[i+1])||
((cbLogicValue[i]==cbLogicValue[i+1])&&(cbCardData[i]<cbCardData[i+1])))
{
//交换位置
cbTempData=cbCardData[i]
cbCardData[i]=cbCardData[i+1]
cbCardData[i+1]=cbTempData
cbTempData=cbLogicValue[i]
cbLogicValue[i]=cbLogicValue[i+1]
cbLogicValue[i+1]=cbTempData
bSorted=false
}
}
bLast--
} while(bSorted==false)

return
}

//混乱扑克
void CGameLogic::RandCardList(BYTE cbCardBuffer[], BYTE cbBufferCount)
{
//CopyMemory(cbCardBuffer,m_cbCardListData,cbBufferCount)
//return
//混乱准备
BYTE cbCardData[CountArray(m_cbCardListData)]
CopyMemory(cbCardData,m_cbCardListData,sizeof(m_cbCardListData))

//混乱扑克
BYTE bRandCount=0,bPosition=0
do
{
//获取随机值，用于解决随机值重复问题  added by hty
int r1=(int)(rand()+time(NULL)+GetTickCount())
srand(r1)
int r=r1+(r1<<3)+(r1>>3)+rand()

bPosition=r%(CountArray(m_cbCardListData)-bRandCount)
cbCardBuffer[bRandCount++]=cbCardData[bPosition]
cbCardData[bPosition]=cbCardData[CountArray(m_cbCardListData)-bRandCount]
} while (bRandCount<cbBufferCount)

return
}

//逻辑数值
func (g *GameLogic) GetCardLogicValue(cardData int) int {
	//cardColor := GetCardColor(cardData)
	cardValue := GetCardValue(cardData)
	if cardValue >0 {
		cardValue = 10
	}
	return  cardValue
}

//对比扑克
bool CGameLogic::CompareCard(BYTE cbFirstData[], BYTE cbNextData[], BYTE cbCardCount,BOOL FirstOX,BOOL NextOX)
{
if(FirstOX!=NextOX)return (FirstOX>NextOX)

if((GetCardType(cbFirstData,cbCardCount) == OX_FIVEKING)&&(GetCardType(cbNextData,cbCardCount)!=OX_FIVEKING)) return true
if((GetCardType(cbFirstData,cbCardCount) != OX_FIVEKING)&&(GetCardType(cbNextData,cbCardCount)==OX_FIVEKING)) return false
//比较牛大小
if(FirstOX==TRUE)
{
//获取点数
BYTE cbNextType=0/*GetCardType(cbNextData,cbCardCount)*/
BYTE cbFirstType=0/*GetCardType(cbFirstData,cbCardCount)*/
///////////11111111111111111111
BYTE cbValue=GetCardLogicValue(cbNextData[3])
cbValue+=GetCardLogicValue(cbNextData[4])

bool firstking=false
bool nextking=false

bool firstDa=false
bool nextDa=false
if(cbValue>10)
{														//防止原值变小
if( (cbNextData[3]==0x4E||cbNextData[4]==0x4F||cbNextData[4]==0x4E||cbNextData[3]==0x4F) /*&& (cbValue<20)*/)
{
//if (cbNextData[4]==0x4F||cbNextData[3]==0x4F)
BYTE ShengYu=0
cbValue=0
for(BYTE i=3i<5i++)//3 4
{
cbValue+=GetCardLogicValue(cbNextData[i])
}
ShengYu=cbValue%10
if (ShengYu>0 )
{
nextDa=true//nextDa是判断4,5有没有利用大王的
}

cbValue=10
}
else
cbValue-=10 //2.3
}
cbNextType=cbValue

//liu
////////////
BYTE bKingCount=0
for(BYTE i=0i<3i++)
{
if (cbNextData[i]==0x4E||cbNextData[i]==0x4F)
{
//if (cbNextData[i]==0x4F)
/*{
	nextDa=true //不用了  nextDa是判断4,5有没有利用大王的
}*/
bKingCount++
}
}
if (bKingCount>0)
{
cbValue=0

BYTE ShengYu=0
for(BYTE i=0i<3i++)//0 1 2
{
cbValue+=GetCardLogicValue(cbNextData[i])
}
ShengYu=cbValue%10
if (ShengYu>0 )
{
nextking=true//是判断1,2,3有没有利用大王的
}
}
/////////////////////////////////
cbValue=0
cbValue=GetCardLogicValue(cbFirstData[3])
cbValue+=GetCardLogicValue(cbFirstData[4])

if(cbValue>10)
{														//防止原值变小
if( (cbFirstData[3]==0x4E||cbFirstData[4]==0x4F||cbFirstData[4]==0x4E||cbFirstData[3]==0x4F) /*&& (cbValue<20)*/)
{
BYTE ShengYu=0
cbValue=0
for(BYTE i=3i<5i++)//0 1 2
{
cbValue+=GetCardLogicValue(cbFirstData[i])
}
ShengYu=cbValue%10
if (ShengYu>0 )
{
firstDa=true
}

cbValue=10
}
else
cbValue-=10 //2.3
}
cbFirstType=cbValue
////////////////
bKingCount=0
for(BYTE i=0i<3i++)
{
if (cbFirstData[i]==0x4E||cbFirstData[i]==0x4F)
{
//if (cbFirstData[i]==0x4F)
/*{
	firstDa=true
}*/
bKingCount++
}
}
if (bKingCount>0)
{
cbValue=0

BYTE ShengYu=0
for(BYTE i=0i<3i++)//0 1 2
{
cbValue+=GetCardLogicValue(cbFirstData[i])
}
ShengYu=cbValue%10
if (ShengYu>0 )
{
firstking=true
}
}
///////////11111111111111111111



if (cbFirstType==cbNextType /*&& cbFirstType==10*/)
{
//同点数大王>小王>...
BYTE cbFirstKing=10
BYTE cbNextKing=10
for (int i=0i<5i++)
{
if (cbFirstData[i]==0x4E)
cbFirstKing=11
else if (cbFirstData[i]==0x4F)
cbFirstKing=12

if (cbNextData[i]==0x4E)
cbNextKing=11
else if (cbNextData[i]==0x4F)
cbNextKing=12
}

if (cbNextKing!=cbFirstKing)
{
return cbFirstKing>cbNextKing
}

/*if (firstking&&nextking)//都采用了大小王
{
	if (firstDa)
	{
		return true
	}
	else if (nextDa)
	{
		return false
	}
}

if (firstking)
{
	return false
}
else if (nextking)
{
	return true
}*/
if ((firstking||firstDa)&&(nextking||nextDa))//都采用了大小王
{

}
else if (firstking||firstDa)
{
return true
}
else if (nextking||nextDa)
{
return false
}
}
//点数判断
if (cbFirstType!=cbNextType) return (cbFirstType>cbNextType)

//switch(cbNextType)
//{
//case OX_FOUR_SAME:		//炸弹牌型
//	{
//		//排序大小
//		BYTE bFirstTemp[MAX_COUNT],bNextTemp[MAX_COUNT]
//		CopyMemory(bFirstTemp,cbFirstData,cbCardCount)
//		CopyMemory(bNextTemp,cbNextData,cbCardCount)
//		SortCardList(bFirstTemp,cbCardCount)
//		SortCardList(bNextTemp,cbCardCount)

//		return GetCardValue(bFirstTemp[MAX_COUNT/2])>GetCardValue(bNextTemp[MAX_COUNT/2])

//		break
//	}
//case OX_THREE_SAME:		//小牛
//	{
//		//排序大小
//		BYTE bFirstTemp[MAX_COUNT],bNextTemp[MAX_COUNT]
//		CopyMemory(bFirstTemp,cbFirstData,cbCardCount)
//		CopyMemory(bNextTemp,cbNextData,cbCardCount)
//		SortCardList(bFirstTemp,cbCardCount)
//		SortCardList(bNextTemp,cbCardCount)

//		return GetCardValue(bFirstTemp[0])>GetCardValue(bNextTemp[0])

//		break
//	}
//}
}

//排序大小
BYTE bFirstTemp[MAX_COUNT],bNextTemp[MAX_COUNT]
CopyMemory(bFirstTemp,cbFirstData,cbCardCount)
CopyMemory(bNextTemp,cbNextData,cbCardCount)
SortCardList(bFirstTemp,cbCardCount)
SortCardList(bNextTemp,cbCardCount)

//比较数值
BYTE cbNextMaxValue=GetCardValue(bNextTemp[0])
BYTE cbFirstMaxValue=GetCardValue(bFirstTemp[0])
if(cbNextMaxValue!=cbFirstMaxValue)return cbFirstMaxValue>cbNextMaxValue

//比较颜色
return GetCardColor(bFirstTemp[0]) > GetCardColor(bNextTemp[0])

return false
}

//获取数值
func (g *GameLogic) GetCardValue(cardData int) int {
	return cardData & LOGIC_MASK_VALUE
}
//获取花色
func (g *GameLogic) GetCardColor(cardData int) int {
	return cardData & LOGIC_MASK_COLOR
}

//////////////////////////////////////////////////////////////////////////
