package config

import (
	"fmt"
	"log"
)

// migrationColumn describes a single ADD-COLUMN migration step.
type migrationColumn struct {
	table    string
	column   string
	alterSQL string
}

// userSignMigrations lists every new column that must be present on the
// user_sign table.  Each entry is applied exactly once: the column is added
// only when it does not already exist in the live schema.
var userSignMigrations = []migrationColumn{
	{
		table:  "user_sign",
		column: "notify_mode",
		alterSQL: "ALTER TABLE user_sign ADD COLUMN notify_mode INT NOT NULL DEFAULT 0" +
			" COMMENT '签到通知模式 0-全部通知 1-仅失败通知 2-仅成功通知'" +
			" AFTER user_number",
	},
	{
		table:    "user_sign",
		column:   "ap_remind",
		alterSQL: "ALTER TABLE user_sign ADD COLUMN ap_remind INT NOT NULL DEFAULT 0 COMMENT '是否开启理智提醒 0-关闭 1-开启' AFTER notify_mode",
	},
	{
		table:    "user_sign",
		column:   "ap_threshold",
		alterSQL: "ALTER TABLE user_sign ADD COLUMN ap_threshold INT NOT NULL DEFAULT 80 COMMENT '理智提醒阈值百分比' AFTER ap_remind",
	},
	{
		table:    "user_sign",
		column:   "ap_notified",
		alterSQL: "ALTER TABLE user_sign ADD COLUMN ap_notified INT NOT NULL DEFAULT 0 COMMENT '理智提醒是否已通知 0-未通知 1-已通知' AFTER ap_threshold",
	},
}

// MigrateDB applies all pending schema migrations to the connected database.
// Each migration step is idempotent: a column is only added when it is not
// already present in the live schema.  If the target table does not exist
// (e.g. fresh deployment where arknights.sql was already applied with the
// new schema, or a deployment that hasn't created the table yet), the
// migration for that table is silently skipped.
func MigrateDB() error {
	if DBEngine == nil {
		return fmt.Errorf("migrate: database engine is not initialized")
	}
	for _, m := range userSignMigrations {
		exists, err := tableExists(m.table)
		if err != nil {
			return err
		}
		if !exists {
			log.Printf("migrate: table %s does not exist, skipping column migrations", m.table)
			break // all remaining entries target the same table
		}
		if err := applyColumnIfMissing(m); err != nil {
			return err
		}
	}

	// Create user_ap_remind table and migrate data from user_sign if needed.
	if err := migrateApRemindTable(); err != nil {
		return err
	}
	if err := migrateGuestSpamTables(); err != nil {
		return err
	}

	return nil
}

// tableExists returns true when the named table exists in the current database.
func tableExists(table string) (bool, error) {
	var count int64
	res := DBEngine.Raw(
		`SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES
		 WHERE TABLE_SCHEMA = DATABASE()
		   AND TABLE_NAME   = ?`,
		table,
	).Scan(&count)
	if res.Error != nil {
		return false, fmt.Errorf("migrate: checking table %s existence: %w", table, res.Error)
	}
	return count > 0, nil
}

// applyColumnIfMissing adds m.column to m.table when it does not already exist.
// It queries INFORMATION_SCHEMA.COLUMNS so the check works on any MySQL-compatible
// database without relying on database-specific DDL extensions.
func applyColumnIfMissing(m migrationColumn) error {
	var count int64
	res := DBEngine.Raw(
		`SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		 WHERE TABLE_SCHEMA = DATABASE()
		   AND TABLE_NAME   = ?
		   AND COLUMN_NAME  = ?`,
		m.table, m.column,
	).Scan(&count)
	if res.Error != nil {
		return fmt.Errorf("migrate: checking column %s.%s: %w", m.table, m.column, res.Error)
	}
	if count > 0 {
		return nil // already present
	}
	if err := DBEngine.Exec(m.alterSQL).Error; err != nil {
		return fmt.Errorf("migrate: adding column %s.%s: %w", m.table, m.column, err)
	}
	log.Printf("migrate: added column %s.%s", m.table, m.column)
	return nil
}

// migrateApRemindTable creates the user_ap_remind table if it does not exist,
// and migrates existing AP reminder data from user_sign (for upgrades from
// the old schema where AP fields were embedded in user_sign).
func migrateApRemindTable() error {
	exists, err := tableExists("user_ap_remind")
	if err != nil {
		return err
	}
	if !exists {
		createSQL := `CREATE TABLE user_ap_remind (
			id varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			user_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '用户名称',
			user_number bigint NULL DEFAULT NULL COMMENT '用户ID',
			ap_threshold int NULL DEFAULT 80 COMMENT '理智提醒阈值百分比',
			ap_notified int NULL DEFAULT 0 COMMENT '理智提醒是否已通知 0-未通知 1-已通知',
			remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			create_time timestamp(0) NULL DEFAULT NULL,
			update_time timestamp(0) NULL DEFAULT NULL,
			PRIMARY KEY (id) USING BTREE
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '理智提醒用户' ROW_FORMAT = Dynamic`
		if err := DBEngine.Exec(createSQL).Error; err != nil {
			return fmt.Errorf("migrate: creating user_ap_remind table: %w", err)
		}
		log.Println("migrate: created table user_ap_remind")

		// Migrate existing data from user_sign.ap_* columns if the old columns exist.
		oldColExists, err := columnExists("user_sign", "ap_remind")
		if err != nil {
			return err
		}
		if oldColExists {
			migrateSQL := `INSERT INTO user_ap_remind (id, user_name, user_number, ap_threshold, ap_notified, create_time, update_time)
				SELECT REPLACE(UUID(), '-', ''), user_name, user_number, ap_threshold, ap_notified, create_time, update_time
				FROM user_sign WHERE ap_remind = 1`
			result := DBEngine.Exec(migrateSQL)
			if result.Error != nil {
				return fmt.Errorf("migrate: migrating AP data from user_sign: %w", result.Error)
			}
			if result.RowsAffected > 0 {
				log.Printf("migrate: migrated %d AP reminder records from user_sign to user_ap_remind", result.RowsAffected)
			}
		}
	}
	return nil
}

// columnExists returns true when the named column exists in the given table.
func columnExists(table, column string) (bool, error) {
	var count int64
	res := DBEngine.Raw(
		`SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		 WHERE TABLE_SCHEMA = DATABASE()
		   AND TABLE_NAME   = ?
		   AND COLUMN_NAME  = ?`,
		table, column,
	).Scan(&count)
	if res.Error != nil {
		return false, fmt.Errorf("migrate: checking column %s.%s: %w", table, column, res.Error)
	}
	return count > 0, nil
}

func migrateGuestSpamTables() error {
	tables := map[string]string{
		"guest_spam_member_risk": `CREATE TABLE guest_spam_member_risk (
			id varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			chat_id bigint NULL DEFAULT NULL,
			user_id bigint NULL DEFAULT NULL,
			user_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			first_seen_at timestamp(0) NULL DEFAULT NULL,
			last_message_at timestamp(0) NULL DEFAULT NULL,
			recent_message_count bigint NULL DEFAULT 0,
			warning_count bigint NULL DEFAULT 0,
			mute_level bigint NULL DEFAULT 0,
			last_penalty_at timestamp(0) NULL DEFAULT NULL,
			remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			create_time timestamp(0) NULL DEFAULT NULL,
			update_time timestamp(0) NULL DEFAULT NULL,
			PRIMARY KEY (id) USING BTREE,
			INDEX idx_guest_spam_member_chat_user (chat_id, user_id)
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam成员风控状态' ROW_FORMAT = Dynamic`,
		"guest_spam_member_activity": `CREATE TABLE guest_spam_member_activity (
			id varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			chat_id bigint NULL DEFAULT NULL,
			user_id bigint NULL DEFAULT NULL,
			activity_day date NOT NULL,
			message_count bigint NULL DEFAULT 0,
			remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			create_time timestamp(0) NULL DEFAULT NULL,
			update_time timestamp(0) NULL DEFAULT NULL,
			PRIMARY KEY (id) USING BTREE,
			INDEX idx_guest_spam_activity_chat_user_day (chat_id, user_id, activity_day),
			INDEX idx_guest_spam_activity_day (activity_day)
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam成员日活跃' ROW_FORMAT = Dynamic`,
		"guest_spam_bot_blacklist": `CREATE TABLE guest_spam_bot_blacklist (
			id varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			bot_id bigint NULL DEFAULT NULL,
			bot_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			bot_user_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			source varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			first_chat_id bigint NULL DEFAULT NULL,
			first_message_id int NULL DEFAULT NULL,
			first_caller_user_id bigint NULL DEFAULT NULL,
			first_caller_chat_id bigint NULL DEFAULT NULL,
			remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			create_time timestamp(0) NULL DEFAULT NULL,
			update_time timestamp(0) NULL DEFAULT NULL,
			PRIMARY KEY (id) USING BTREE,
			UNIQUE INDEX uk_guest_spam_bot_id (bot_id)
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam bot全局黑名单' ROW_FORMAT = Dynamic`,
		"guest_spam_log": `CREATE TABLE guest_spam_log (
			id varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
			chat_id bigint NULL DEFAULT NULL,
			chat_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			message_id int NULL DEFAULT NULL,
			guest_bot_id bigint NULL DEFAULT NULL,
			guest_bot_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			guest_bot_user varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			caller_user_id bigint NULL DEFAULT NULL,
			caller_user_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			caller_chat_id bigint NULL DEFAULT NULL,
			caller_chat_name varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			action varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			reason varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			detail varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			remark varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
			create_time timestamp(0) NULL DEFAULT NULL,
			update_time timestamp(0) NULL DEFAULT NULL,
			PRIMARY KEY (id) USING BTREE,
			INDEX idx_guest_spam_log_chat (chat_id),
			INDEX idx_guest_spam_log_bot (guest_bot_id),
			INDEX idx_guest_spam_log_action (action)
		) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = 'guest spam处理日志' ROW_FORMAT = Dynamic`,
	}
	for table, createSQL := range tables {
		exists, err := tableExists(table)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if err := DBEngine.Exec(createSQL).Error; err != nil {
			return fmt.Errorf("migrate: creating %s table: %w", table, err)
		}
		log.Printf("migrate: created table %s", table)
	}
	return nil
}
