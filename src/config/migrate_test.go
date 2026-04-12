package config

import (
	"strings"
	"testing"
)

// TestUserSignMigrationsCount verifies that the expected four AP-reminder
// columns are declared in the migration list.
func TestUserSignMigrationsCount(t *testing.T) {
	want := 4
	if got := len(userSignMigrations); got != want {
		t.Errorf("want %d migration entries, got %d", want, got)
	}
}

// TestUserSignMigrationsColumns verifies that every expected column has its
// own migration entry with the correct table name.
func TestUserSignMigrationsColumns(t *testing.T) {
	wantColumns := []string{"notify_mode", "ap_remind", "ap_threshold", "ap_notified"}
	wantTable := "user_sign"

	columnSet := make(map[string]bool)
	for _, m := range userSignMigrations {
		columnSet[m.column] = true
	}

	for _, col := range wantColumns {
		if !columnSet[col] {
			t.Errorf("migration entry missing for column %q", col)
		}
	}

	for _, m := range userSignMigrations {
		if m.table != wantTable {
			t.Errorf("migration for column %q targets table %q, want %q", m.column, m.table, wantTable)
		}
	}
}

// TestUserSignMigrationsAlterSQL verifies that each ALTER SQL string is
// syntactically reasonable: it must reference ALTER TABLE, the correct table
// name, ADD COLUMN, and the column name it is registering.
func TestUserSignMigrationsAlterSQL(t *testing.T) {
	for _, m := range userSignMigrations {
		upper := strings.ToUpper(m.alterSQL)

		if !strings.Contains(upper, "ALTER TABLE") {
			t.Errorf("column %q: alterSQL does not contain ALTER TABLE: %q", m.column, m.alterSQL)
		}
		if !strings.Contains(upper, strings.ToUpper(m.table)) {
			t.Errorf("column %q: alterSQL does not reference table %q: %q", m.column, m.table, m.alterSQL)
		}
		if !strings.Contains(upper, "ADD COLUMN") {
			t.Errorf("column %q: alterSQL does not contain ADD COLUMN: %q", m.column, m.alterSQL)
		}
		if !strings.Contains(upper, strings.ToUpper(m.column)) {
			t.Errorf("column %q: alterSQL does not reference column name: %q", m.column, m.alterSQL)
		}
		if !strings.Contains(upper, "AFTER") {
			t.Errorf("column %q: alterSQL does not contain AFTER clause for column positioning: %q", m.column, m.alterSQL)
		}
	}
}

// TestUserSignMigrationsDefaults verifies that each column's ALTER SQL
// specifies a DEFAULT value.
func TestUserSignMigrationsDefaults(t *testing.T) {
	wantDefaults := map[string]string{
		"notify_mode":  "DEFAULT 0",
		"ap_remind":    "DEFAULT 0",
		"ap_threshold": "DEFAULT 80",
		"ap_notified":  "DEFAULT 0",
	}

	for _, m := range userSignMigrations {
		want, ok := wantDefaults[m.column]
		if !ok {
			continue
		}
		if !strings.Contains(strings.ToUpper(m.alterSQL), want) {
			t.Errorf("column %q: expected DEFAULT clause %q in alterSQL: %q", m.column, want, m.alterSQL)
		}
	}
}

// TestMigrateDBNilEngine verifies that MigrateDB returns a descriptive error
// when DBEngine has not been initialized.
func TestMigrateDBNilEngine(t *testing.T) {
	// Save and restore the global DBEngine.
	original := DBEngine
	defer func() { DBEngine = original }()

	DBEngine = nil

	err := MigrateDB()
	if err == nil {
		t.Fatal("want error when DBEngine is nil, got nil")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("want 'not initialized' in error message, got: %v", err)
	}
}

// TestMigrationColumnNoDuplicates verifies that each column name appears
// exactly once in the migration list.
func TestMigrationColumnNoDuplicates(t *testing.T) {
	seen := make(map[string]int)
	for _, m := range userSignMigrations {
		seen[m.column]++
	}
	for col, count := range seen {
		if count > 1 {
			t.Errorf("column %q appears %d times in migrations, want 1", col, count)
		}
	}
}

// TestUserSignMigrationsColumnOrder verifies that the AFTER clauses create
// the correct column ordering: notify_mode → ap_remind → ap_threshold → ap_notified,
// matching the order defined in arknights.sql and user_sign.go.
func TestUserSignMigrationsColumnOrder(t *testing.T) {
	wantAfter := map[string]string{
		"notify_mode":  "AFTER USER_NUMBER",
		"ap_remind":    "AFTER NOTIFY_MODE",
		"ap_threshold": "AFTER AP_REMIND",
		"ap_notified":  "AFTER AP_THRESHOLD",
	}

	for _, m := range userSignMigrations {
		want, ok := wantAfter[m.column]
		if !ok {
			continue
		}
		if !strings.Contains(strings.ToUpper(m.alterSQL), want) {
			t.Errorf("column %q: expected %q in alterSQL, got: %q", m.column, want, m.alterSQL)
		}
	}
}

// TestMigrateDBAllSameTable verifies that all migration entries target
// the same table, which the table-existence check relies on to break early.
func TestMigrateDBAllSameTable(t *testing.T) {
	if len(userSignMigrations) == 0 {
		t.Fatal("no migrations defined")
	}
	table := userSignMigrations[0].table
	for _, m := range userSignMigrations[1:] {
		if m.table != table {
			t.Errorf("migration for %q targets table %q, want %q (same as first entry)", m.column, m.table, table)
		}
	}
}
