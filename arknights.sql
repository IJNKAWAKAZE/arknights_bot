SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for group_invite
-- ----------------------------
DROP TABLE IF EXISTS `group_invite`;
CREATE TABLE `group_invite`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `group_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '群组名称',
  `group_number` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '邀请人名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '邀请人ID',
  `member_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '被邀请人名称',
  `member_number` bigint NULL DEFAULT NULL COMMENT '被邀请人ID',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群组邀请记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_joined
-- ----------------------------
DROP TABLE IF EXISTS `group_joined`;
CREATE TABLE `group_joined`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '0',
  `group_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '群组名称',
  `group_number` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `news` int NOT NULL DEFAULT 0 COMMENT '消息推送开关0关闭1开启',
  `reg` int NULL DEFAULT NULL COMMENT '群规消息ID',
  `welcome` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '入群欢迎信息',
  `birthday` int NULL DEFAULT 0 COMMENT '干员生日推送开关0关闭1开启',
  `request_mode` int NULL DEFAULT NULL COMMENT '群验证模式0正常加入1申请加入',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '机器人加入群组记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_lottery
-- ----------------------------
DROP TABLE IF EXISTS `group_lottery`;
CREATE TABLE `group_lottery`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `group_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '群组名称',
  `group_number` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `status` int NULL DEFAULT 1 COMMENT '抽奖状态0关闭1开启2暂停报名',
  `end_time` timestamp(0) NULL DEFAULT NULL COMMENT '抽奖报名结束时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群组抽奖信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for group_lottery_detail
-- ----------------------------
DROP TABLE IF EXISTS `group_lottery_detail`;
CREATE TABLE `group_lottery_detail`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `lottery_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '群组抽奖ID',
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `lottery_number` int NULL DEFAULT NULL COMMENT '用户选择号码',
  `status` int NULL DEFAULT 0 COMMENT '中奖状态0未中奖1已中奖',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '群组抽奖详情' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_account
-- ----------------------------
DROP TABLE IF EXISTS `user_account`;
CREATE TABLE `user_account`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `hypergryph_token` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '鹰角账号token',
  `skland_token` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '森空岛token',
  `skland_cred` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '森空岛cred',
  `skland_id` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '森空岛ID',
  `server_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '服务器名称',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户账户信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_gacha
-- ----------------------------
DROP TABLE IF EXISTS `user_gacha`;
CREATE TABLE `user_gacha`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `uid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色UID',
  `pool_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '卡池名称',
  `pool_order` int NULL DEFAULT NULL COMMENT '卡池顺序',
  `char_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色名称',
  `is_new` tinyint(1) NULL DEFAULT NULL COMMENT '是否NEW标',
  `rarity` int NULL DEFAULT NULL COMMENT '星级',
  `ts` bigint NULL DEFAULT NULL COMMENT '抽卡时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户抽卡信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_player
-- ----------------------------
DROP TABLE IF EXISTS `user_player`;
CREATE TABLE `user_player`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `account_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '账号表ID',
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `uid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色UID',
  `server_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '服务器名称',
  `player_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色名称',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户绑定角色信息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_sign
-- ----------------------------
DROP TABLE IF EXISTS `user_sign`;
CREATE TABLE `user_sign`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `notify_mode` int NULL DEFAULT 0 COMMENT '签到通知模式 0-全部通知 1-仅失败通知 2-仅成功通知',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '自动签到用户' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_ap_remind
-- ----------------------------
DROP TABLE IF EXISTS `user_ap_remind`;
CREATE TABLE `user_ap_remind`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `user_number` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `ap_threshold` int NULL DEFAULT 80 COMMENT '理智提醒阈值百分比',
  `ap_notified` int NULL DEFAULT 0 COMMENT '理智提醒是否已通知 0-未通知 1-已通知',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '理智提醒用户' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for guest_spam_member_risk
-- ----------------------------
DROP TABLE IF EXISTS `guest_spam_member_risk`;
CREATE TABLE `guest_spam_member_risk`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `chat_id` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
  `first_seen_at` timestamp(0) NULL DEFAULT NULL COMMENT '首次见到时间',
  `last_message_at` timestamp(0) NULL DEFAULT NULL COMMENT '最近普通发言时间',
  `recent_message_count` bigint NULL DEFAULT 0 COMMENT '近期普通发言数',
  `warning_count` bigint NULL DEFAULT 0 COMMENT 'guest spam警告数',
  `mute_level` bigint NULL DEFAULT 0 COMMENT '禁言阶梯',
  `last_penalty_at` timestamp(0) NULL DEFAULT NULL COMMENT '最近处罚时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_guest_spam_member_chat_user`(`chat_id`, `user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam成员风控状态' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for guest_spam_member_activity
-- ----------------------------
DROP TABLE IF EXISTS `guest_spam_member_activity`;
CREATE TABLE `guest_spam_member_activity`  (
  `id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `chat_id` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `user_id` bigint NULL DEFAULT NULL COMMENT '用户ID',
  `day` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '日期yyyyMMdd',
  `message_count` bigint NULL DEFAULT 0 COMMENT '当天普通发言数',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_guest_spam_activity_chat_user_day`(`chat_id`, `user_id`, `day`) USING BTREE,
  INDEX `idx_guest_spam_activity_day`(`day`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam成员日活跃' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for guest_spam_bot_blacklist
-- ----------------------------
DROP TABLE IF EXISTS `guest_spam_bot_blacklist`;
CREATE TABLE `guest_spam_bot_blacklist`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `bot_id` bigint NULL DEFAULT NULL COMMENT 'guest bot ID',
  `bot_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'guest bot名称',
  `bot_user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'guest bot用户名',
  `source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '拉黑来源',
  `first_chat_id` bigint NULL DEFAULT NULL COMMENT '首次判定群组',
  `first_message_id` int NULL DEFAULT NULL COMMENT '首次判定消息',
  `first_caller_user_id` bigint NULL DEFAULT NULL COMMENT '首次触发用户',
  `first_caller_chat_id` bigint NULL DEFAULT NULL COMMENT '首次触发频道/身份',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_guest_spam_bot_id`(`bot_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam bot全局黑名单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for guest_spam_log
-- ----------------------------
DROP TABLE IF EXISTS `guest_spam_log`;
CREATE TABLE `guest_spam_log`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `chat_id` bigint NULL DEFAULT NULL COMMENT '群组ID',
  `chat_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '群组名称',
  `message_id` int NULL DEFAULT NULL COMMENT '消息ID',
  `guest_bot_id` bigint NULL DEFAULT NULL COMMENT 'guest bot ID',
  `guest_bot_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'guest bot名称',
  `guest_bot_user` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'guest bot用户名',
  `caller_user_id` bigint NULL DEFAULT NULL COMMENT '触发用户ID',
  `caller_user_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '触发用户名称',
  `caller_chat_id` bigint NULL DEFAULT NULL COMMENT '触发频道/身份ID',
  `caller_chat_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '触发频道/身份名称',
  `action` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '处理动作',
  `reason` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '处理原因',
  `detail` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '详情',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `create_time` timestamp(0) NULL DEFAULT NULL,
  `update_time` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_guest_spam_log_chat`(`chat_id`) USING BTREE,
  INDEX `idx_guest_spam_log_bot`(`guest_bot_id`) USING BTREE,
  INDEX `idx_guest_spam_log_action`(`action`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam处理日志' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
