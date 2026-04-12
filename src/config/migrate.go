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
