package database

import "github.com/forbole/bdjuno/v4/database/utils"

func (db *Db) SaveDenom(txHash, denomID, name, schema, sender, uri string) error {
	_, err := db.SQL.Exec(
		`INSERT INTO nft_denom (transaction_hash, id, name, schema, sender, uri) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
		txHash, denomID, name, schema, sender, uri)
	return err
}

func (db *Db) UpdateDenom(denomID, owner string) error {
	_, err := db.SQL.Exec(`UPDATE nft_denom SET owner = $1 WHERE id = $2`, owner, denomID)
	return err
}

func (tx *DbTx) SaveNFT(txHash string, tokenID uint64, denomID, name, description, uri string, tags []string, sender, recipient string) error {
	_, err := tx.Exec(`INSERT INTO nft_nft (transaction_hash, id, denom_id, name, description, uri, tags, sender, recipient, uniq_id) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (id, denom_id) DO UPDATE SET uniq_id = EXCLUDED.uniq_id`, txHash, tokenID, denomID, name, description, uri, tags,
		sender, recipient, utils.FormatUniqID(tokenID, denomID))
	return err
}

func (db *Db) UpdateNFT(id, denomID, name, uri, dataJSON, dataText string) error {
	_, err := db.SQL.Exec(`UPDATE nft_nft SET name = $1, uri = $2, data_json = $3, data_text = $4 WHERE id = $5 AND denom_id = $6`, name, uri, dataJSON, dataText, id, denomID)
	return err
}

func (tx *DbTx) UpdateNFTOwner(id, denomID, owner string) error {
	_, err := tx.Exec(`UPDATE nft_nft SET owner = $1 WHERE id = $2 AND denom_id = $3`, owner, id, denomID)
	return err
}

func (tx *DbTx) BurnNFT(id, denomID string) error {
	_, err := tx.Exec(`UPDATE nft_nft SET burned = true WHERE id = $1 AND denom_id = $2`, id, denomID)
	return err
}

func (tx *DbTx) UpdateNFTHistory(txHash string, tokenID uint64, denomID, from, to string, timestamp uint64) error {
	_, err := tx.Exec(`INSERT INTO nft_transfer_history (transaction_hash, id, denom_id, old_owner, new_owner, timestamp, uniq_id) 
		VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING`,
		txHash, tokenID, denomID, from, to, timestamp, utils.FormatUniqID(tokenID, denomID))
	return err
}
