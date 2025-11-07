package service

const (
	findHistoryQuery                   = `SELECT id, old_value, new_value, status, action, created_date, created_by FROM %s WHERE %s = ? ORDER BY id desc limit ?`
	insertHistoryQuery                 = `INSERT INTO %s ( %s, old_value, new_value, status, action, created_by, created_date, store_id ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	createMigrationFlagQuery           = `INSERT INTO migration_flag ( table_name, action, last_update_id, ` + "`limit`" + `) VALUES (?, ?, ?, ?)`
	findLastUpdateIdMigrationFlagQuery = `SELECT * FROM migration_flag WHERE table_name = ? ORDER BY id desc limit 1`
)
