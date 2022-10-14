package accounts_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/alteredAccount"
	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/stretchr/testify/require"
)

func generateApiAccounts(n int) []*alteredAccount.AlteredAccount {
	ret := make([]*alteredAccount.AlteredAccount, n)

	for i := 0; i < n; i++ {
		ret[i] = generateApiAccount()
	}

	return ret
}

func generateApiAccount() *alteredAccount.AlteredAccount {
	return &alteredAccount.AlteredAccount{
		Address: "erd1q",
		Balance: testscommon.GenerateRandomBigInt().String(),
		Nonce:   rand.Uint64(),
	}
}

func TestAlteredAccountsProcessor_ProcessAccounts(t *testing.T) {
	t.Parallel()

	ap := accounts.NewAlteredAccountsProcessor()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		apiAccounts := generateApiAccounts(10)
		alteredAccounts, err := ap.ProcessAccounts(apiAccounts)
		require.Nil(t, err)
		requireAlteredAccountsProcessedSuccessfully(t, apiAccounts, alteredAccounts)
	})

	t.Run("nil api account, should skip it", func(t *testing.T) {
		t.Parallel()

		apiAccounts := generateApiAccounts(10)
		apiAccounts[0] = nil
		alteredAccounts, err := ap.ProcessAccounts(apiAccounts)
		require.Nil(t, err)
		requireAlteredAccountsProcessedSuccessfully(t, apiAccounts[1:], alteredAccounts)
	})

	t.Run("invalid balance, should return error", func(t *testing.T) {
		t.Parallel()

		apiAccounts := generateApiAccounts(10)
		apiAccounts[4].Balance = "balance"
		alteredAccounts, err := ap.ProcessAccounts(apiAccounts)
		require.Nil(t, alteredAccounts)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "balance"))
	})
}

func requireAlteredAccountsProcessedSuccessfully(
	t *testing.T,
	apiAccounts []*alteredAccount.AlteredAccount,
	processedAccounts []*schemaV2.AccountBalanceUpdate,
) {
	require.Equal(t, len(apiAccounts), len(processedAccounts))

	for idx := range apiAccounts {
		requireAlteredAccountProcessedSuccessfully(t, apiAccounts[idx], processedAccounts[idx])
	}
}

func requireAlteredAccountProcessedSuccessfully(
	t *testing.T,
	apiAccount *alteredAccount.AlteredAccount,
	processedAccount *schemaV2.AccountBalanceUpdate,
) {
	balance, err := utility.GetBigIntBytesFromStr(apiAccount.Balance)
	require.Nil(t, err)

	expectedAccount := &schemaV2.AccountBalanceUpdate{
		Address: []byte(apiAccount.Address),
		Balance: balance,
		Nonce:   int64(apiAccount.Nonce),
	}
	require.Equal(t, expectedAccount, processedAccount)
}
