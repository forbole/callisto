package types

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
)

// DidDoc represents the x/did doc
type DidDoc struct {
	ID                   string
	Context              []string
	Controller           []string
	VerificationMethod   []*didtypes.VerificationMethod
	Authentication       []string
	AssertionMethod      []string
	CapabilityInvocation []string
	CapabilityDelegation []string
	KeyAgreement         []string
	Service              []*didtypes.Service
	AlsoKnownAs          []string
	VersionId            string
	FromAddress          string
	Height               int64
}

// NewDidDoc allows to build a new DidDoc instance
func NewDidDoc(id string,
	context []string,
	controller []string,
	verificationMethod []*didtypes.VerificationMethod,
	authentication []string,
	assertionMethod []string,
	capabilityInvocation []string,
	capabilityDelegation []string,
	keyAgreement []string,
	service []*didtypes.Service,
	alsoKnownAs []string,
	versionId string,
	fromAddress string,
	height int64) *DidDoc {
	return &DidDoc{
		ID:                   id,
		Context:              context,
		Controller:           controller,
		VerificationMethod:   verificationMethod,
		Authentication:       authentication,
		AssertionMethod:      assertionMethod,
		CapabilityInvocation: capabilityInvocation,
		CapabilityDelegation: capabilityDelegation,
		KeyAgreement:         keyAgreement,
		Service:              service,
		AlsoKnownAs:          alsoKnownAs,
		VersionId:            versionId,
		FromAddress:          fromAddress,
		Height:               height,
	}
}
