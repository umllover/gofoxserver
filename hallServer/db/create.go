package db

var userUpdateSql = [][]string{
	0: []string{
		`SET FOREIGN_KEY_CHECKS=0;`,

		`CREATE TABLE accountsmember (
		UserID bigint(11) NOT NULL COMMENT '用户标识',
		MemberOrder tinyint(4) NOT NULL DEFAULT '0' COMMENT '会员标识',
		UserRight int(11) NOT NULL COMMENT '用户权限',
		MemberOverDate timestamp NULL DEFAULT NULL COMMENT '会员期限',
		PRIMARY KEY (UserID),
		KEY IX_OverDate (MemberOverDate),
		KEY IX_UserID (UserID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE create_room_info (
		user_id bigint(11) NOT NULL COMMENT '用户索引',
		room_name varchar(255) NOT NULL,
		kind_id int(11) NOT NULL COMMENT '房间索引',
		service_id int(11) NOT NULL COMMENT '游戏标识',
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '录入日期',
		node_id int(11) NOT NULL DEFAULT '0' COMMENT '在哪个服务器上',
		room_id int(11) NOT NULL DEFAULT '0' COMMENT '房间id',
		num int(11) NOT NULL DEFAULT '0' COMMENT '局数',
		status int(11) NOT NULL DEFAULT '0',
		Public int(11) NOT NULL DEFAULT '0' COMMENT '是否公开',
		max_player_cnt int(11) NOT NULL DEFAULT '0' COMMENT '最多几个玩家进入',
		pay_type int(11) NOT NULL DEFAULT '1' COMMENT '支付方式 1是全服 2是AA',
		other_info varchar(255) DEFAULT NULL COMMENT '其他配置 json格式',
		PRIMARY KEY (room_id),
		KEY IX_GameScoreLocker_UserID_ServerID (user_id,service_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE gamescoreinfo (
		UserID bigint(11) NOT NULL COMMENT '用户 ID',
		Score bigint(20) NOT NULL COMMENT '用户积分（货币）',
		Revenue bigint(20) NOT NULL COMMENT '游戏税收',
		InsureScore bigint(20) NOT NULL COMMENT '银行金币',
		WinCount int(11) NOT NULL COMMENT '胜局数目',
		LostCount int(11) NOT NULL COMMENT '输局数目',
		DrawCount int(11) NOT NULL COMMENT '和局数目',
		FleeCount int(11) NOT NULL COMMENT '逃局数目',
		AllLogonTimes int(11) NOT NULL COMMENT '总登陆次数',
		PlayTimeCount int(11) NOT NULL COMMENT '游戏时间',
		OnLineTimeCount int(11) NOT NULL COMMENT '在线时间',
		LastLogonIP varchar(15) NOT NULL COMMENT '上次登陆 IP',
		LastLogonDate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '上次登陆时间',
		LastLogonMachine varchar(32) NOT NULL COMMENT '登录机器',
		RegisterIP varchar(15) NOT NULL COMMENT '注册 IP',
		RegisterMachine varchar(32) NOT NULL COMMENT '注册机器',
		PRIMARY KEY (UserID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE gamescorelocker (
		UserID bigint(11) NOT NULL COMMENT '用户索引',
		KindID int(11) NOT NULL COMMENT '房间索引',
		ServerID int(11) NOT NULL COMMENT '游戏标识',
		HallNodeID int(11) NOT NULL,
		GameNodeID int(11) DEFAULT NULL COMMENT '在哪个服务器上',
		roomid int(11) NOT NULL COMMENT '进出索引',
		EnterIP varchar(255) NOT NULL COMMENT '进入地址',
		EnterMachine varchar(32) NOT NULL COMMENT '进入机器',
		CollectDate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '录入日期',
		PRIMARY KEY (UserID),
		KEY IX_GameScoreLocker_UserID_ServerID (UserID,ServerID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE goods_live (
		id int(11) NOT NULL COMMENT '物品id',
		left_amount int(10) NOT NULL COMMENT '剩余的数量',
		trade_time int(11) NOT NULL COMMENT '交易次数',
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE inc_userid (
		node_id int(11) NOT NULL,
		inc_id bigint(11) NOT NULL,
		PRIMARY KEY (node_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE mail (
		mail_id int(11) NOT NULL,
		user_id int(11) NOT NULL COMMENT '邮件id',
		mail_type int(11) NOT NULL,
		context varchar(255) NOT NULL,
		creator_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		sender varchar(255) NOT NULL,
		title varchar(255) NOT NULL,
		PRIMARY KEY (mail_id,user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE race_msg_info (
		MsgID int(11) NOT NULL AUTO_INCREMENT,
		SendTimes int(11) DEFAULT NULL COMMENT '还需要发送多少次，发完删除该记录',
		IntervalTime int(11) DEFAULT NULL,
		Context text,
		MsgType int(11) DEFAULT NULL,
		PRIMARY KEY (MsgID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE room_id (
		id int(11) NOT NULL,
		node_id int(11) NOT NULL,
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE room_record (
		room_id int(11) NOT NULL,
		kind_id int(11) NOT NULL,
		user_id bigint(11) NOT NULL COMMENT '创建房间的玩家id',
		status int(11) NOT NULL COMMENT '游戏状态 ',
		room_name varchar(255) NOT NULL COMMENT '房间名字',
		jion_user varchar(255) DEFAULT NULL COMMENT '进入的玩家id',
		PRIMARY KEY (room_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE token_record (
		room_id int(11) NOT NULL,
		user_id bigint(11) NOT NULL,
		tokenType int(11) NOT NULL,
		amount int(11) NOT NULL,
		status int(11) NOT NULL,
		creator_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		KindID int(11) NOT NULL,
		ServerId int(11) NOT NULL,
		PRIMARY KEY (room_id,user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;`,

		`CREATE TABLE userattr (
		UserID bigint(11) NOT NULL,
		UnderWrite varchar(255) DEFAULT NULL COMMENT '个性签名',
		FaceID smallint(6) NOT NULL COMMENT '头像标识',
		CustomID int(11) NOT NULL DEFAULT '0' COMMENT '自定标识',
		UserMedal int(11) NOT NULL DEFAULT '0' COMMENT '用户奖牌',
		Experience int(11) NOT NULL COMMENT '经验数值',
		LoveLiness int(11) NOT NULL COMMENT '用户魅力',
		UserRight int(11) NOT NULL DEFAULT '0' COMMENT '用户权限',
		MasterRight int(11) NOT NULL COMMENT '管理权限',
		MasterOrder tinyint(11) NOT NULL DEFAULT '0' COMMENT '管理等级',
		PlayTimeCount int(11) NOT NULL DEFAULT '0' COMMENT '游戏时间',
		OnLineTimeCount int(11) NOT NULL COMMENT '在线时间',
		HeadImgUrl varchar(255) NOT NULL COMMENT '头像',
		Gender tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别',
		NickName varchar(255) NOT NULL,
		elect_uid bigint(11) NOT NULL COMMENT '推举人id',
		star int(11) NOT NULL COMMENT '点赞数',
		Sign varchar(11) NOT NULL COMMENT '个性签名',
		phome_number varchar(11) NOT NULL COMMENT '电话号码',
		PRIMARY KEY (UserID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE userextrainfo (
		UserId bigint(11) NOT NULL,
		MbPayTotal int(11) NOT NULL COMMENT '手机充值总额',
		MbVipLevel int(11) NOT NULL COMMENT '手机VIP等级',
		PayMbVipUpgrade int(11) NOT NULL COMMENT '手机VIP升级，所需充值数（vip最高级时该值为0）',
		MbTicket int(11) NOT NULL COMMENT '手机兑换券数量',
		PRIMARY KEY (UserId)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE usertoken (
		UserID bigint(11) NOT NULL,
		Currency int(11) NOT NULL COMMENT '游戏豆',
		RoomCard int(11) NOT NULL COMMENT '房卡数',
		PRIMARY KEY (UserID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_day_times (
		user_id bigint(11) NOT NULL,
		key_id int(11) NOT NULL,
		v int(11) NOT NULL,
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id,key_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_mask_code (
		user_id bigint(11) NOT NULL,
		phome_number varchar(11) NOT NULL COMMENT '电话号码',
		mask_code int(11) NOT NULL COMMENT '验证按',
		creator_time varchar(255) DEFAULT NULL,
		PRIMARY KEY (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_offline_handler (
		id int(11) NOT NULL AUTO_INCREMENT,
		user_id bigint(11) NOT NULL,
		h_type varchar(255) NOT NULL,
		context varchar(255) NOT NULL,
		expiry_time timestamp NULL DEFAULT NULL,
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_spread (
		user_id bigint(11) NOT NULL,
		spread_uid bigint(11) NOT NULL COMMENT '被我领取的推广人id',
		status int(11) NOT NULL DEFAULT '0' COMMENT '是否领取了这个人的奖励',
		PRIMARY KEY (user_id,spread_uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_times (
		user_id bigint(11) NOT NULL,
		key_id int(11) NOT NULL,
		v int(11) NOT NULL,
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id,key_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE user_week_times (
		user_id bigint(11) NOT NULL,
		key_id int(11) NOT NULL,
		v int(11) NOT NULL,
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id,key_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	},
	1: []string{
		"ALTER TABLE create_room_info add user_cnt int(11) NOT NULL DEFAULT 0 COMMENT '加入的玩家数';",
	},

	2: []string{
		"ALTER TABLE token_record add play_cnt int(11) NOT NULL DEFAULT 0  COMMENT '可玩的局数';",
	},

	3: []string{
		`CREATE TABLE record_outcard_ddz (
		RecordID bigint(11) NOT NULL COMMENT '记录ID',
		CreateTime int(11) NOT NULL COMMENT '创建时间',
		CardData text COMMENT '牌数据，数组转成字符串',
		PRIMARY KEY (RecordID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	},

	4: []string{
		`CREATE TABLE record_outcard_ddz_king (
		RecordID bigint(11) NOT NULL COMMENT '记录ID八王表',
		CreateTime int(11) NOT NULL COMMENT '创建时间',
		CardData text COMMENT '牌数据，数组转成字符串',
		PRIMARY KEY (RecordID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	},

	5: []string{
		"ALTER TABLE record_outcard_ddz MODIFY RecordID BIGINT(11) NOT NULL AUTO_INCREMENT;",
		"ALTER TABLE record_outcard_ddz_king MODIFY RecordID BIGINT(11) NOT NULL AUTO_INCREMENT;",
	},
}

///////////////////////////////////////////////////// log db /////////////////////////////////////////////////
var statsUpdateSql = [][]string{
	0: []string{
		`SET FOREIGN_KEY_CHECKS=0;`,

		`CREATE TABLE activity (
		activity_name varchar(11) NOT NULL COMMENT '活动名',
		activity_type int(11) NOT NULL COMMENT '活动类别',
		activity_begin timestamp NULL DEFAULT NULL COMMENT '活动开始时间',
		activity_end timestamp NULL DEFAULT NULL COMMENT '活动开始时间',
		PRIMARY KEY (activity_name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE consum_log (
		recode_id int(11) NOT NULL AUTO_INCREMENT,
		user_id bigint(11) NOT NULL COMMENT '用户索引',
		consum_type int(11) NOT NULL DEFAULT '0' COMMENT '消费类型 0钻石 1开房 3道具',
		consum_num int(11) NOT NULL DEFAULT '0' COMMENT '消费数量',
		consum_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '消费时间',
		PRIMARY KEY (recode_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE draw_award_log (
		id int(11) NOT NULL COMMENT '活动id。 和程序保持一致',
		draw_id int(11) NOT NULL COMMENT '领取奖励的key',
		description varchar(255) NOT NULL COMMENT '活动描述',
		draw_count bigint(11) NOT NULL COMMENT '活动可以领取的次数',
		draw_type int(11) NOT NULL COMMENT '领取类型，1是永久，2是每日领取，3是每周领取 ',
		amount int(11) NOT NULL COMMENT '奖励数量',
		item_type int(11) NOT NULL COMMENT '领取的物品类型， 1是钻石，',
		draw_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '领取奖励的时间',
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE getin_room_log (
		recode_id int(11) NOT NULL AUTO_INCREMENT COMMENT '加入游戏数据记录的Id',
		room_id int(11) NOT NULL COMMENT '房间id',
		user_id bigint(11) NOT NULL COMMENT '用户索引',
		kind_id int(11) NOT NULL COMMENT '房间索引',
		service_id int(11) NOT NULL COMMENT '游戏标识',
		room_name varchar(255) NOT NULL,
		node_id int(11) NOT NULL DEFAULT '0' COMMENT '在哪个服务器上',
		num int(11) NOT NULL DEFAULT '0' COMMENT '局数',
		status int(11) NOT NULL DEFAULT '0',
		public int(11) NOT NULL DEFAULT '0' COMMENT '公房加入 0否 1是',
		max_player_cnt int(11) NOT NULL DEFAULT '0' COMMENT '最多几个玩家进入',
		pay_type int(11) NOT NULL COMMENT '支付方式 1是全服 2是AA',
		type_getIn int(11) DEFAULT NULL COMMENT '加入房间类型 0列表加入 2输房号加入 3快速加入 4点击链接加入',
		getIn_time timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '进入房间时间',
		PRIMARY KEY (recode_id)
		) ENGINE=InnoDB AUTO_INCREMENT=439 DEFAULT CHARSET=utf8;`,

		`CREATE TABLE globalspreadinfo (
		ID int(11) NOT NULL,
		RegisterGrantScore int(11) NOT NULL COMMENT '注册时赠送金币数目',
		PlayTimeCount int(11) NOT NULL COMMENT '在线时长（单位：秒）',
		PlayTimeGrantScore int(11) NOT NULL COMMENT '根据在线时长赠送金币数目',
		FillGrantRate decimal(18,2) NOT NULL COMMENT '充值赠送比率',
		BalanceRate decimal(18,2) NOT NULL COMMENT '结算赠送比率',
		MinBalanceScore int(11) NOT NULL COMMENT '结算最小值',
		PRIMARY KEY (ID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE mall_buy_log (
		goods_id int(15) NOT NULL,
		rmb int(11) NOT NULL,
		diamond int(11) NOT NULL,
		name varchar(255) NOT NULL COMMENT '商品名称',
		left_cnt int(11) NOT NULL COMMENT '剩余数量',
		special_offer int(11) NOT NULL COMMENT '特价',
		give_present int(11) NOT NULL COMMENT '赠送',
		special_offer_begin timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '特价开始时间',
		special_offer_end timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '特价结束时间',
		goods_type varchar(255) DEFAULT NULL COMMENT '类别',
		buy_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '购买时间',
		PRIMARY KEY (goods_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE recharge_log (
		OnLineID int(11) NOT NULL COMMENT '订单标识',
		PayAmount int(11) NOT NULL COMMENT '实付金额',
		UserID bigint(11) NOT NULL COMMENT '用户标识',
		PayType varchar(255) NOT NULL DEFAULT '' COMMENT '支付类型',
		GoodsID int(11) NOT NULL COMMENT '物品id',
		RechangeTime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '冲值时间',
		PRIMARY KEY (OnLineID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE recommend_log (
		sub_elect_uid bigint(11) NOT NULL COMMENT '被推举人人id',
		elect_uid bigint(11) NOT NULL COMMENT '推举人id',
		elect_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '领取时间',
		PRIMARY KEY (sub_elect_uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE room_log (
		recode_id int(11) NOT NULL AUTO_INCREMENT COMMENT '房间数据记录的Id',
		room_id int(11) NOT NULL COMMENT '房间id',
		user_id bigint(11) NOT NULL COMMENT '用户索引',
		room_name varchar(255) NOT NULL,
		kind_id int(11) NOT NULL COMMENT '房间索引',
		service_id int(11) NOT NULL COMMENT '游戏标识',
		node_id int(11) NOT NULL COMMENT '在哪个服务器上',
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '录入日期',
		end_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '结束日期',
		create_others int(11) NOT NULL COMMENT '是否为他人开房 0否，1是',
		pay_type int(11) NOT NULL COMMENT '支付方式 1是全服 2是AA',
		timeout_nostart int(11) NOT NULL COMMENT '是否超时未开始游戏 0否  1是 ',
		start_endError int(11) NOT NULL COMMENT '是否开始但非正常解散房间 0 否 1是',
		nomal_open int(11) NOT NULL COMMENT '是否正常开房 0否 1是',
		PRIMARY KEY (recode_id)
		) ENGINE=InnoDB AUTO_INCREMENT=153 DEFAULT CHARSET=utf8;`,

		`CREATE TABLE systemgrantcount (
		DateID int(11) NOT NULL,
		RegisterIP char(15) NOT NULL COMMENT '注册地址',
		RegisterMachine varchar(32) NOT NULL COMMENT '注册机器',
		GrantScore bigint(20) NOT NULL COMMENT '赠送金币',
		GrantCount bigint(20) NOT NULL COMMENT '赠送次数',
		CollectDate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '收集时间',
		PRIMARY KEY (DateID,RegisterIP),
							KEY IX_SystemGrantCount_CollectDate (CollectDate),
							KEY IX_SystemGrantCount_RegisterIP (RegisterIP),
							KEY IX_SystemGrantCount_RegisterMachine (RegisterMachine)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE systemstatusinfo (
		StatusName varchar(32) NOT NULL COMMENT '状态名字',
		StatusValue int(11) NOT NULL COMMENT '状态数值',
		StatusString text NOT NULL COMMENT '状态字符',
		StatusTip varchar(50) NOT NULL COMMENT '状态显示名称',
		StatusDescription varchar(100) NOT NULL COMMENT '字符的描述',
		SortID int(11) NOT NULL,
		PRIMARY KEY (StatusName)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`CREATE TABLE systemstreaminfo (
		DateID int(11) NOT NULL COMMENT '日期标识',
		WebLogonSuccess int(11) NOT NULL COMMENT '登录成功',
		WebRegisterSuccess int(11) NOT NULL COMMENT '注册成功',
		GameLogonSuccess int(11) NOT NULL COMMENT '登录成功',
		GameRegisterSuccess int(11) NOT NULL COMMENT '注册成功',
		CollectDate timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '录入时间',
		PRIMARY KEY (DateID),
		KEY IX_SystemStreamInfo_CollectDate (CollectDate)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	},
	1: []string{
		`ALTER TABLE room_log
		 MODIFY COLUMN create_time  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '录入日期' AFTER node_id,
		 MODIFY COLUMN end_time  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束日期' AFTER create_time,
		 CHANGE COLUMN timeout_nostart game_end_type  int(11) NOT NULL COMMENT '游戏结束类型 0是常规结束 1是游戏解散 2是玩家请求解散 3是没开始就解散' AFTER pay_type,
		 CHANGE COLUMN start_endError room_end_type  int(11) NOT NULL COMMENT '解散房间类型 1出错解散房间 2正常解散房间' AFTER game_end_type;`,
	},
	2: []string{
		`ALTER TABLE consum_log
		MODIFY COLUMN recode_id  int(11) NOT NULL AUTO_INCREMENT FIRST ;`,
	},
}
