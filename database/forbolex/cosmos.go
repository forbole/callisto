package forbolex

import "github.com/desmos-labs/juno/types"

// SaveTx implements Db.
// Performs a no-op since we don't care about transaction.
func (db *Db) SaveTx(tx *types.Tx) error {
	return nil
}

// HasValidator implements Db.
// Performs a no-op since we don't care about validators.
func (db *Db) HasValidator(address string) (bool, error) {
	return true, nil
}

// SaveValidator implements Db.
// Performs a no-op since we don't care about validators.
func (db *Db) SaveValidator(address, publicKey string) error {
	return nil
}

// SaveCommitSig implements Db.
// Performs a no-op since we don't care about commits.
func (db *Db) SaveCommitSig(commitSig *types.CommitSig) error {
	return nil
}

// SaveMessage implements Db.
// Performs a no-op since we don't care about messages.
func (db *Db) SaveMessage(msg *types.Message) error {
	return nil
}
