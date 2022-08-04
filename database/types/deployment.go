package types

// DeploymentParamsRow represents a single row inside the deployment_params table
type DeploymentParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewDeploymentParamsRow builds a new DeploymentParamsRow instance
func NewDeploymentParamsRow(
	params string, height int64,
) DeploymentParamsRow {
	return DeploymentParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m DeploymentParamsRow) Equal(n DeploymentParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
