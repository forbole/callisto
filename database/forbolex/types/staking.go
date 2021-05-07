package types

// ValidatorInfoRow represents a single row inside the "validator_info" table
type ValidatorInfoRow struct {
	ValidatorOperAddr         string `db:"val_oper_addr"`
	ValidatorSelfDelegateAddr string `db:"val_self_delegate_addr"`
}
