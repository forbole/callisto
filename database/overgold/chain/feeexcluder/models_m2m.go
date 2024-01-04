package feeexcluder

import (
	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK M2MTariffFees

// toM2MTariffFeesDatabase - mapping func to a database model.
func toM2MTariffFeesDatabase(tariffID, feesID uint64) db.FeeExcluderM2MTariffFees {
	return db.FeeExcluderM2MTariffFees{
		TariffID: tariffID,
		FeesID:   feesID,
	}
}

// BLOCK M2MTariffTariffs

// toM2MTariffFeesDatabase - mapping func to a database model.
func toM2MTariffTariffsDatabase(tariffID, tariffsID uint64) db.FeeExcluderM2MTariffTariffs {
	return db.FeeExcluderM2MTariffTariffs{
		TariffID:  tariffID,
		TariffsID: tariffsID,
	}
}

// BLOCK M2MGenesisStateTariffs

// toM2MTariffFeesDatabase - mapping func to a database model.
func toM2MGenesisStateTariffs(genesisStateID, tariffsID uint64) db.FeeExcluderM2MGenesisStateTariffs {
	return db.FeeExcluderM2MGenesisStateTariffs{
		GenesisStateID: genesisStateID,
		TariffsID:      tariffsID,
	}
}
