package accounts_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/stretchr/testify/require"
)

func createNotarizedBlocks() []*api.NotarizedBlock {
	alteredAcc1 := &outport.AlteredAccount{
		Nonce:   4,
		Address: "erd1a",
		Balance: "444",
		Tokens:  nil,
	}
	alteredAcc2 := &outport.AlteredAccount{
		Nonce:   4,
		Address: "erd1b",
		Balance: "444",
		Tokens: []*outport.AccountTokenData{
			{
				Nonce:      1,
				Identifier: "identifier1",
				Balance:    "111",
				Properties: "properties1",
				MetaData:   nil,
			},
			{
				Nonce:      2,
				Identifier: "identifier2",
				Balance:    "222",
				Properties: "properties2",
				MetaData:   nil,
			},
		},
	}
	alteredAcc3 := &outport.AlteredAccount{
		Nonce:   7,
		Address: "erd1c",
		Balance: "444",
		Tokens: []*outport.AccountTokenData{
			{
				Nonce:      5,
				Identifier: "identifier3",
				Balance:    "555",
				Properties: "properties3",
				MetaData: &esdt.MetaData{
					Nonce:      6,
					Name:       []byte("name"),
					Creator:    []byte("creator"),
					Royalties:  666,
					Hash:       []byte("hash"),
					URIs:       [][]byte{[]byte("uri1"), []byte("uri2")},
					Attributes: []byte("attributes"),
				},
			},
		},
	}

	return []*api.NotarizedBlock{
		{
			Shard:           1,
			AlteredAccounts: []*outport.AlteredAccount{alteredAcc1},
		},
		{
			Shard:           2,
			AlteredAccounts: []*outport.AlteredAccount{alteredAcc2, alteredAcc3},
		},
	}
}

func TestAlteredAccountsProcessor_ProcessAccounts(t *testing.T) {
	t.Parallel()

	ap := accounts.NewAlteredAccountsProcessor()

	processedAcc1 := &schema.AccountBalanceUpdate{
		Address: []byte("erd1a"),
		Balance: big.NewInt(444).Bytes(),
		Nonce:   4,
		Tokens:  nil,
	}
	processedAcc2 := &schema.AccountBalanceUpdate{
		Address: []byte("erd1b"),
		Balance: big.NewInt(444).Bytes(),
		Nonce:   4,
		Tokens: []*schema.AccountTokenData{
			{
				Nonce:      1,
				Identifier: "identifier1",
				Balance:    big.NewInt(111).Bytes(),
				Properties: "properties1",
				MetaData:   nil,
			},
			{
				Nonce:      2,
				Identifier: "identifier2",
				Balance:    big.NewInt(222).Bytes(),
				Properties: "properties2",
				MetaData:   nil,
			},
		},
	}
	processedAcc3 := &schema.AccountBalanceUpdate{
		Address: []byte("erd1c"),
		Balance: big.NewInt(444).Bytes(),
		Nonce:   7,
		Tokens: []*schema.AccountTokenData{
			{
				Nonce:      5,
				Identifier: "identifier3",
				Balance:    big.NewInt(555).Bytes(),
				Properties: "properties3",
				MetaData: &schema.MetaData{
					Nonce:      6,
					Name:       []byte("name"),
					Creator:    []byte("creator"),
					Royalties:  666,
					Hash:       []byte("hash"),
					URIs:       [][]byte{[]byte("uri1"), []byte("uri2")},
					Attributes: []byte("attributes"),
				},
			},
		},
	}

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()

		expectedAccounts := []*schema.AccountBalanceUpdate{processedAcc1, processedAcc2, processedAcc3}
		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Equal(t, expectedAccounts, res)
		require.Nil(t, err)
	})

	t.Run("nil notarized block, should skip it", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()
		notarizedBlocks[0] = nil

		expectedAccounts := []*schema.AccountBalanceUpdate{processedAcc2, processedAcc3}
		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Equal(t, expectedAccounts, res)
		require.Nil(t, err)
	})

	t.Run("nil altered account, should skip it", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()
		notarizedBlocks[1].AlteredAccounts[0] = nil

		expectedAccounts := []*schema.AccountBalanceUpdate{processedAcc1, processedAcc3}
		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Equal(t, expectedAccounts, res)
		require.Nil(t, err)
	})

	t.Run("nil token, should skip it", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()
		notarizedBlocks[1].AlteredAccounts[0].Tokens[1] = nil

		processedAcc2Copy := *processedAcc2
		processedAcc2Copy.Tokens = []*schema.AccountTokenData{processedAcc2.Tokens[0]}
		expectedAccounts := []*schema.AccountBalanceUpdate{processedAcc1, &processedAcc2Copy, processedAcc3}
		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Equal(t, expectedAccounts, res)
		require.Nil(t, err)
	})

	t.Run("invalid native balance, should return error", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()
		notarizedBlocks[1].AlteredAccounts[1].Balance = "invalidNativeBalance"

		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Nil(t, res)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "invalidNativeBalance"))
	})

	t.Run("invalid token balance, should return error", func(t *testing.T) {
		t.Parallel()

		notarizedBlocks := createNotarizedBlocks()
		notarizedBlocks[1].AlteredAccounts[0].Tokens[0].Balance = "invalidTokenBalance"

		res, err := ap.ProcessAccounts(notarizedBlocks)
		require.Nil(t, res)
		require.NotNil(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "invalidTokenBalance"))
	})
}
