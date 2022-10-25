package types

type SoftwareUpgradePlanRow struct {
	ProposalID    uint64 `db:"proposal_id"`
	PlanName      string `db:"plan_name"`
	UpgradeHeight int64  `db:"upgrade_height"`
	Info          string `db:"info"`
	Height        int64  `db:"height"`
}

func NewSoftwareUpgradePlanRow(
	proposalID uint64, planName string, upgradeHeight int64, info string, height int64,
) SoftwareUpgradePlanRow {
	return SoftwareUpgradePlanRow{
		ProposalID:    proposalID,
		PlanName:      planName,
		UpgradeHeight: upgradeHeight,
		Info:          info,
		Height:        height,
	}
}
