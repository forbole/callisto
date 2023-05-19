package database

import "fmt"

// Prune implements db.PruningDb
func (db *Db) Prune(height int64) error {
	// Prune default tables
	err := db.Database.Prune(height)
	if err != nil {
		return fmt.Errorf("error while pruning db: %s", err)
	}

	// Prune modules
	err = db.pruneBank(height)
	if err != nil {
		return fmt.Errorf("error while pruning bank: %s", err)
	}

	err = db.pruneSlashing(height)
	if err != nil {
		return fmt.Errorf("error while pruning slashing: %s", err)
	}

	return nil
}

func (db *Db) pruneBank(height int64) error {
	_, err := db.SQL.Exec(`DELETE FROM supply WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning supply: %s", err)
	}
	return nil
}

func (db *Db) pruneSlashing(height int64) error {
	_, err := db.SQL.Exec(`DELETE FROM validator_signing_info WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning validator signing info: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM slashing_params WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning slashing params: %s", err)
	}

	return nil
}
