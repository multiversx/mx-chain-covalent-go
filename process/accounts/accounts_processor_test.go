package accounts_test

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAccountsProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter)
		expectedErr error
	}{
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return nil, &mock.AccountsAdapterStub{}, &mock.PubKeyConverterStub{}
			},
			expectedErr: covalent.ErrNilShardCoordinator,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, nil, &mock.PubKeyConverterStub{}
			},
			expectedErr: covalent.ErrNilAccountsAdapter,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, &mock.AccountsAdapterStub{}, nil
			},
			expectedErr: covalent.ErrNilPubKeyConverter,
		},
		{
			args: func() (process.ShardCoordinator, covalent.AccountsAdapter, core.PubkeyConverter) {
				return &mock.ShardCoordinatorMock{}, &mock.AccountsAdapterStub{}, &mock.PubKeyConverterStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := accounts.NewAccountsProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestAccountsProcessor_ProcessAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(&mock.ShardCoordinatorMock{}, &mock.AccountsAdapterStub{}, &mock.PubKeyConverterStub{})

	tx1 := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes(),
	}
	tx2 := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes(),
	}
	txs := []*schema.Transaction{tx1, tx2}

	scr := &schema.SCResult{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes(),
	}

	scrs := []*schema.SCResult{scr}

	receipt := &schema.Receipt{
		Sender: testscommon.GenerateRandomBytes(),
	}

	receipts := []*schema.Receipt{receipt}

	ret := ap.ProcessAccounts(txs, scrs, receipts)

	require.Len(t, ret, 7)
}
