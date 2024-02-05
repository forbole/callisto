package iscn

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/types"

	"github.com/cosmos/cosmos-sdk/codec"
	iscntypes "github.com/likecoin/likecoin-chain/v4/x/iscn/types"

	"github.com/rs/zerolog/log"
)

func HandleGenesis(appState map[string]json.RawMessage, cdc codec.Codec, db *database.Db) error {
	log.Debug().Str("module", "iscn").Msg("parsing genesis")

	err := GetGenesisIscnRecords(appState, db, cdc)

	if err != nil {
		return fmt.Errorf("error while storing genesis iscn records: %s", err)
	}

	return nil

}

// GetGenesisIscnRecords parses the given appState and saves the genesis iscn records in db
func GetGenesisIscnRecords(appState map[string]json.RawMessage, db *database.Db, cdc codec.Codec) error {

	var genState iscntypes.GenesisState
	err := cdc.UnmarshalJSON(appState[iscntypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading iscn genesis data: %s", err)
	}

	// Store the records
	iscnRecords := make([]types.Record, len(genState.IscnRecords))
	contentRecords := make([]types.IscnRecord, len(genState.ContentIdRecords))
	var height int64 = 0
	ipld := make([]string, len(genState.IscnRecords))

	// Store iscn_records
	for i, record := range genState.IscnRecords {

		mapIscnRecords := map[string]interface{}{}
		err = json.Unmarshal(record, &mapIscnRecords)
		if err != nil {
			return fmt.Errorf("error when trying to unmarshal record at index %d as JSON: %v", i, err)
		}

		idAny, ok := mapIscnRecords["@id"]
		if !ok {
			return fmt.Errorf("error: couldn't find iscn ID field at index %d", i)
		}
		idStr, ok := idAny.(string)
		if !ok {
			return fmt.Errorf("error: iscn ID at index %d is not in string format", i)
		}
		iscnID, err := iscntypes.ParseIscnId(idStr)
		if err != nil {
			return fmt.Errorf("error: invalid iscn ID at index %d : %v", i, err)
		}

		fingerprints, ok := mapIscnRecords["contentFingerprints"]
		if !ok {
			return fmt.Errorf("error: couldn't find content fingerprints field for iscn record with ID %s", iscnID.String())
		}

		stakeholders, ok := mapIscnRecords["stakeholders"]
		if !ok {
			return fmt.Errorf("error: couldn't find stakeholders field for iscn record with ID %s", iscnID.String())
		}

		contentMetadata, ok := mapIscnRecords["contentMetadata"]
		if !ok {
			return fmt.Errorf("error: couldn't find content metadata field for iscn record with ID %s", iscnID.String())
		}

		ipldRecord, ok := mapIscnRecords["ipld"]
		if !ok {
			return fmt.Errorf("error: couldn't find ipld field for iscn record with ID %s", iscnID.String())
		}

		iscnFingerprints := fingerprints.([]string)
		iscnStakeholders := stakeholders.([]iscntypes.IscnInput)
		iscnContentMetadata := contentMetadata.(iscntypes.IscnInput)

		iscnRecords[i] = types.NewRecord(iscnID.String(), "", iscnFingerprints, iscnStakeholders, iscnContentMetadata)
		ipld[i] = ipldRecord.(string)
	}

	// Store content_id_records
	for index, contentIDRecord := range genState.ContentIdRecords {
		var latestVersion = contentIDRecord.LatestVersion
		ownerAddress := contentIDRecord.Owner

		iscnID, err := iscntypes.ParseIscnId(contentIDRecord.IscnId)
		if err != nil {
			return fmt.Errorf("error: couldn't parse iscn ID %s in content ID record entries: %w", contentIDRecord.IscnId, err)
		}

		id := iscnID.String()

		contentRecords[index] = types.NewIscnRecord(ownerAddress, id, latestVersion, ipld[index], iscnRecords[index], height)

		// Save each record in db
		err = db.SaveIscnRecord(contentRecords[index])
		if err != nil {
			return err
		}
	}

	return nil

}
