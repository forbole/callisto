package database

// GetTotalBlocks implements database.Database
func (db *Db) GetTotalBlocks() (int64, error) {
	var blockCount int64
	err := db.Sql.QueryRow(`SELECT count(*) FROM block;`).Scan(&blockCount)
	return blockCount, err
}
