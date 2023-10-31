// Package cli Copyright (c) 2016-2021 Shanghai Bianjie AI Technology Inc. (licensed under the Apache License, Version 2.0)
// Modifications Copyright (c) 2021-present Crypto.org (licensed under the Apache License, Version 2.0)
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenName        = "name"
	FlagTokenDescription = "description"
	FlagTokenFileURI     = "uri"
	FlagRecipient        = "recipient"
	FlagOwner            = "owner"
	FlagTokenTags        = "tags"

	FlagDenomName = "name"
	FlagDenomID   = "denom-id"
	FlagSchema    = "schema"
	FlagDenomURI  = "uri"
)

var (
	FsIssueDenom  = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditNFT     = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueDenom.String(FlagSchema, "", "Denom data structure definition")
	FsIssueDenom.String(FlagDenomName, "", "The name of the denom")
	FsIssueDenom.String(FlagDenomURI, "", "URI of the denom")

	FsMintNFT.String(FlagTokenFileURI, "", "URI for supplemental off-chain tokenFile data (should return a File object)")
	FsMintNFT.String(FlagRecipient, "", "Receiver of the nft, if not filled, the default is the sender of the transaction")
	FsMintNFT.StringSlice(FlagTokenTags, []string{}, "The search tags of the nft")
	FsMintNFT.String(FlagTokenDescription, "", "The description of the nft")
	FsMintNFT.String(FlagTokenName, "", "The name of the nft")

	FsEditNFT.String(FlagTokenFileURI, "[do-not-modify]", "URI for supplemental off-chain tokenFile data (should return a File object)")
	FsEditNFT.StringSlice(FlagTokenTags, []string{"[do-not-modify]"}, "The search tags of the nft")
	FsEditNFT.String(FlagTokenDescription, "[do-not-modify]", "TThe description of the nft")
	FsEditNFT.String(FlagTokenName, "[do-not-modify]", "The name of the nft")

	FsQuerySupply.String(FlagOwner, "", "The owner of the nft")

	FsQueryOwner.String(FlagDenomID, "", "The name of the collection")
}
