package types

type MarkerParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewMarkerParamsRow builds a new MarkerParamsRow instance
func NewMarkerParamsRow(
	params string, height int64,
) MarkerParamsRow {
	return MarkerParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m MarkerParamsRow) Equal(n MarkerParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}

// --------------------------------------------------------------------------------------------------------------------

// MarkerAccountRow represents a single row inside the marker account table
type MarkerAccountRow struct {
	Address                string `db:"address"`
	AccessControl          string `db:"access_control"`
	AllowGovernanceControl bool   `db:"allow_governance_control"`
	Denom                  string `db:"denom"`
	MarkerType             string `db:"marker_type"`
	Status                 string `db:"status"`
	TotalSupply            string `db:"total_supply"`
	Height                 int64  `db:"height"`
}

// NewMarkerAccountRow allows to easily create a new MarkerAccountRow
func NewMarkerAccountRow(
	address string,
	accessControl string,
	allowGovernanceControl bool,
	denom string,
	markerType string,
	status string,
	totalSupply string,
	height int64,
) MarkerAccountRow {
	return MarkerAccountRow{
		Address:                address,
		AccessControl:          accessControl,
		AllowGovernanceControl: allowGovernanceControl,
		Denom:                  denom,
		MarkerType:             markerType,
		Status:                 status,
		TotalSupply:            totalSupply,
		Height:                 height,
	}
}

// Equals return true if two MarkerAccountRow are the same
func (w MarkerAccountRow) Equals(v MarkerAccountRow) bool {
	return w.Address == v.Address &&
		w.AccessControl == v.AccessControl &&
		w.AllowGovernanceControl == v.AllowGovernanceControl &&
		w.Denom == v.Denom &&
		w.MarkerType == v.MarkerType &&
		w.Status == v.Status &&
		w.TotalSupply == v.TotalSupply &&
		w.Height == v.Height
}
