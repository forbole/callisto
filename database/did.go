package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/lib/pq"
)

// SaveDidDoc allows to store the given did doc inside the database
func (db *Db) SaveDidDoc(didDoc *types.DidDoc) error {
	verMethodBz, err := json.Marshal(&didDoc.VerificationMethod)
	if err != nil {
		return fmt.Errorf("error while marshaling did doc verification method: %s", err)
	}

	serviceBz, err := json.Marshal(&didDoc.Service)
	if err != nil {
		return fmt.Errorf("error while marshaling did doc service: %s", err)
	}

	stmt := `
INSERT INTO did_doc (id, context, controller, verification_method, authentication,
	assertion_method, capability_invocation, capability_delegation, 
	key_agreement, service, also_known_as, version_id, from_address, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (id) DO UPDATE 
    SET context = excluded.context,
        controller = excluded.controller,
        verification_method = excluded.verification_method,
        authentication = excluded.authentication,
        assertion_method = excluded.assertion_method,
        capability_invocation = excluded.capability_invocation,
        capability_delegation = excluded.capability_delegation,
        key_agreement = excluded.key_agreement,
        service = excluded.service,
        also_known_as = excluded.also_known_as,
        version_id = excluded.version_id,
		from_address = excluded.from_address,
        height = excluded.height
WHERE did_doc.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, didDoc.ID, pq.StringArray(didDoc.Context), pq.StringArray(didDoc.Controller), string(verMethodBz),
		pq.StringArray(didDoc.Authentication), pq.StringArray(didDoc.AssertionMethod), pq.StringArray(didDoc.CapabilityInvocation), pq.StringArray(didDoc.CapabilityDelegation),
		pq.StringArray(didDoc.KeyAgreement), string(serviceBz), pq.StringArray(didDoc.AlsoKnownAs), didDoc.VersionID, didDoc.FromAddress, didDoc.Height)
	if err != nil {
		return fmt.Errorf("error while storing did doc: %s", err)
	}

	return nil
}

// DeleteDidDoc removes did doc data from the database
func (db *Db) DeleteDidDoc(id, versionID string) error {
	stmt := `DELETE FROM did_doc WHERE id = $1 AND version_id = $2`
	_, err := db.SQL.Exec(stmt, id, versionID)

	if err != nil {
		return fmt.Errorf("error while deleting did doc: %s", err)
	}
	return nil
}
