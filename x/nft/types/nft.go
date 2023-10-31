// Copyright (c) 2016-2021 Shanghai Bianjie AI Technology Inc. (licensed under the Apache License, Version 2.0)
// Modifications Copyright (c) 2021-present Crypto.org (licensed under the Apache License, Version 2.0)
package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v4/x/nft/exported"
)

var _ exported.NFT = BaseNFT{}

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id, name string, description string, owner sdk.AccAddress, uri string, tags []string) BaseNFT {
	return BaseNFT{
		Id:          id,
		Name:        name,
		Description: description,
		Owner:       owner.String(),
		URI:         uri,
		Tags:        tags,
	}
}

// GetID return the id of BaseNFT
func (bnft BaseNFT) GetID() string {
	return bnft.Id
}

// GetName return the name of BaseNFT
func (bnft BaseNFT) GetName() string {
	return bnft.Name
}

// GetName return the description of BaseNFT
func (bnft BaseNFT) GetDescription() string {
	return bnft.Description
}

// GetOwner return the owner of BaseNFT
func (bnft BaseNFT) GetOwner() sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(bnft.Owner)

	if err != nil {
		panic(fmt.Errorf("couldn't convert %q to account address: %v", bnft.Owner, err))
	}

	return owner
}

// GetURI return the URI of BaseNFT
func (bnft BaseNFT) GetURI() string {
	return bnft.URI
}

// GetData return the Tags of BaseNFT
func (bnft BaseNFT) GetTags() []string {
	return bnft.Tags
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts)
}
