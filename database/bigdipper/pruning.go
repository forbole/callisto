package bigdipper

// Prune implements db.PruningDb
func (db *Db) Prune(height int64) error {
	// Prune default tables
	err := db.Database.Prune(height)
	if err != nil {
		return err
	}

	// Prune modules
	err = db.pruneBank(height)
	if err != nil {
		return err
	}

	err = db.pruneStaking(height)
	if err != nil {
		return err
	}

	err = db.pruneMint(height)
	if err != nil {
		return err
	}

	err = db.pruneDistribution(height)
	if err != nil {
		return err
	}

	err = db.pruneSlashing(height)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) pruneBank(height int64) error {
	_, err := db.Sql.Exec(`DELETE FROM supply WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM account_balance WHERE height = $1`, height)
	return err
}

func (db *Db) pruneStaking(height int64) error {
	_, err := db.Sql.Exec(`DELETE FROM staking_pool WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM validator_commission WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM validator_voting_power WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM validator_status WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM delegation WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM unbonding_delegation WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM redelegation WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM double_sign_vote WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM double_sign_evidence WHERE height = $1`, height)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) pruneMint(height int64) error {
	_, err := db.Sql.Exec(`DELETE FROM inflation WHERE height = $1`, height)
	return err
}

func (db *Db) pruneDistribution(height int64) error {
	_, err := db.Sql.Exec(`DELETE FROM community_pool WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM validator_commission_amount WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM delegation_reward WHERE height = $1`, height)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) pruneSlashing(height int64) error {
	_, err := db.Sql.Exec(`DELETE FROM validator_signing_info WHERE height = $1`, height)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`DELETE FROM slashing_params WHERE height = $1`, height)
	if err != nil {
		return err
	}

	return nil
}
