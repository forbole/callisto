package gov_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/forbole/bdjuno/v4/modules/gov"
)

func TestGetProposalMedata(t *testing.T) {
	testCases := []struct {
		name        string
		metadata    string
		expMetadata string
	}{
		{
			name:        "valid text metadata",
			metadata:    "This is a text metadata",
			expMetadata: "This is a text metadata",
		},
		{
			name:        "valid URL metadata - text",
			metadata:    "https://ipfs.desmos.network/ipfs/QmfPWhiVFCWFxaEd18NBz59TLy8UyT8biB1Gx7B67v3XW8",
			expMetadata: "This is a text proposal",
		},
		{
			name:        "valid URL metadata - non text (image)",
			metadata:    "https://ipfs.desmos.network/ipfs/QmWGu6Egvyydohb3pu12Q2iJpmxAXDNUvSnRQhtNncoH3p",
			expMetadata: "https://ipfs.desmos.network/ipfs/QmWGu6Egvyydohb3pu12Q2iJpmxAXDNUvSnRQhtNncoH3p",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			metadata, err := gov.GetProposalMetadata(tc.metadata)
			require.NoError(t, err)
			require.Equal(t, tc.expMetadata, metadata)
		})
	}
}
