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

	err = db.pruneStaking(height)
	if err != nil {
		return fmt.Errorf("error while pruning staking: %s", err)
	}

	err = db.pruneMint(height)
	if err != nil {
		return fmt.Errorf("error while pruning mint: %s", err)
	}

	err = db.pruneDistribution(height)
	if err != nil {
		return fmt.Errorf("error while pruning distribution: %s", err)
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

func (db *Db) pruneStaking(height int64) error {
	_, err := db.SQL.Exec(`DELETE FROM staking_pool WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning staking pool: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM validator_commission WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning validator commission: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM validator_voting_power WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning validator voting power: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM validator_status WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning validator status: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM double_sign_vote WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning double sign votes: %s", err)
	}

	_, err = db.SQL.Exec(`DELETE FROM double_sign_evidence WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning double sign evidence: %s", err)
	}

	return nil
}

func (db *Db) pruneMint(height int64) error {
	_, err := db.SQL.Exec(`DELETE FROM inflation WHERE height = $1`, height)
	return fmt.Errorf("error while pruning inflation: %s", err)
}

func (db *Db) pruneDistribution(height int64) error {
	_, err := db.SQL.Exec(`DELETE FROM community_pool WHERE height = $1`, height)
	if err != nil {
		return fmt.Errorf("error while pruning community pool: %s", err)
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
