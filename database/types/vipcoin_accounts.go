package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	extratypes "git.ooo.ua/vipcoin/chain/x/types"
	"github.com/lib/pq"
)

type (
	// DBAccount represents a single row inside the "vipcoin_chain_accounts_accounts" table
	DBAccount struct {
		Address    string         `db:"address"`
		Hash       string         `db:"hash"`
		PublicKey  string         `db:"public_key"`
		Kinds      pq.Int32Array  `db:"kinds"`
		State      int32          `db:"state"`
		Extras     ExtraDB        `db:"extras"`
		Affiliates pq.Int64Array  `db:"affiliates"`
		Wallets    pq.StringArray `db:"wallets"`
	}

	// DBAffiliates represents a single row inside the "vipcoin_chain_accounts_affiliates" table
	DBAffiliates struct {
		Id              uint64  `db:"id"`
		Address         string  `db:"address"`
		AffiliationKind int32   `db:"affiliation_kind"`
		Extras          ExtraDB `db:"extras"`
	}

	// DBRegisterUser represents a single row inside the "vipcoin_chain_accounts_register_user" table
	DBRegisterUser struct {
		Creator               string  `db:"creator"`
		Address               string  `db:"address"`
		Hash                  string  `db:"hash"`
		PublicKey             string  `db:"public_key"`
		HolderWallet          string  `db:"holder_wallet"`
		RefRewardWallet       string  `db:"ref_reward_wallet"`
		HolderWalletExtras    ExtraDB `db:"holder_wallet_extras"`
		RefRewardWalletExtras ExtraDB `db:"ref_reward_wallet_extras"`
		ReferrerHash          string  `db:"referrer_hash"`
	}

	// DBSetKinds represents a single row inside the "vipcoin_chain_accounts_set_kinds" table
	DBSetKinds struct {
		Creator string        `db:"creator"`
		Hash    string        `db:"hash"`
		Kinds   pq.Int32Array `db:"kinds"`
	}

	// DBSetAffiliateAddress represents a single row inside the "vipcoin_chain_accounts_set_affiliate_address" table
	DBSetAffiliateAddress struct {
		Creator    string `db:"creator"`
		Hash       string `db:"hash"`
		OldAddress string `db:"old_address"`
		NewAddress string `db:"new_address"`
	}

	// DBAccountMigrate represents a single row inside the "vipcoin_chain_accounts_account_migrate" table
	DBAccountMigrate struct {
		Creator   string `db:"creator"`
		Address   string `db:"address"`
		Hash      string `db:"hash"`
		PublicKey string `db:"public_key"`
	}

	// DBSetAccountExtra represents a single row inside the "vipcoin_chain_accounts_set_extra" table
	DBSetAccountExtra struct {
		Creator string  `db:"creator"`
		Hash    string  `db:"hash"`
		Extras  ExtraDB `db:"extras"`
	}

	// DBSetAffiliateExtra represents a single row inside the "vipcoin_chain_accounts_set_affiliate_extra" table
	DBSetAffiliateExtra struct {
		Creator         string  `db:"creator"`
		AccountHash     string  `db:"account_hash"`
		AffiliationHash string  `db:"affiliation_hash"`
		Extras          ExtraDB `db:"extras"`
	}

	// DBSetState represents a single row inside the "vipcoin_chain_accounts_set_state" table
	DBSetState struct {
		Creator string `db:"creator"`
		Hash    string `db:"hash"`
		State   int32  `db:"state"`
		Reason  string `db:"reason"`
	}

	// DBCreateAccount represents a single row inside the "vipcoin_chain_accounts_create_account" table
	DBCreateAccount struct {
		Creator   string        `db:"creator"`
		Hash      string        `db:"hash"`
		Address   string        `db:"address"`
		PublicKey string        `db:"public_key"`
		Kinds     pq.Int32Array `db:"kinds"`
		State     int32         `db:"state"`
		Extras    ExtraDB       `db:"extras"`
	}

	// DBAddAffiliate represents a single row inside the "vipcoin_chain_accounts_add_affiliate" table
	DBAddAffiliate struct {
		Creator         string  `db:"creator"`
		AccountHash     string  `db:"account_hash"`
		AffiliationHash string  `db:"affiliation_hash"`
		Affiliation     int32   `db:"affiliation"`
		Extras          ExtraDB `db:"extras"`
	}

	// ExtraDB helprs type
	ExtraDB struct {
		Extras []extratypes.Extra
	}
)

// Make the ExtraDB struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (e ExtraDB) Value() (driver.Value, error) {
	return json.Marshal(e.Extras)
}

// Make the ExtraDB struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (e *ExtraDB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e.Extras)
}
