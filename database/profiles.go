package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v2/types"
)

func (db *Db) SaveProfilesParams(params types.ProfilesParams) error {
	nicknameParamsBz, err := json.Marshal(&params.NicknameParams)
	if err != nil {
		return fmt.Errorf("error while marshaling Nickname params: %s", err)
	}
	dTagParamsBz, err := json.Marshal(&params.DTagParams)
	if err != nil {
		return fmt.Errorf("error while marshaling DTag params: %s", err)
	}
	bioParamsBz, err := json.Marshal(&params.BioParams)
	if err != nil {
		return fmt.Errorf("error while marshaling Bio params: %s", err)
	}
	oracleParamsBz, err := json.Marshal(&params.OracleParams)
	if err != nil {
		return fmt.Errorf("error while marshaling Oracle params: %s", err)
	}

	stmt := `
INSERT INTO profiles_params (nickname_params, d_tag_params, bio_params, oracle_params, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET nickname_params = excluded.nickname_params,
		d_tag_params = excluded.d_tag_params,
		bio_params = excluded.bio_params,
		oracle_params = excluded.oracle_params,
      	height = excluded.height
WHERE profiles_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt,
		string(nicknameParamsBz),
		string(dTagParamsBz),
		string(bioParamsBz),
		string(oracleParamsBz),
		params.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}
