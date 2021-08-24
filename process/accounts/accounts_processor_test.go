package accounts_test

import (
	"errors"
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/accounts"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
	"math/big"
	"strconv"
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

func TestAccountsProcessor_ProcessAccounts_InvalidUserAccountHandler_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{
			LoadAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
				return nil, nil
			}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_InvalidLoadAccount_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{
			LoadAccountCalled: func(address []byte) (vmcommon.AccountHandler, error) {
				return nil, errors.New("load account error")
			}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_NotInSameShard_ExpectZeroAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{SelfID: 4},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountStub{}},
		&mock.PubKeyConverterStub{})

	tx := &schema.Transaction{
		Receiver: testscommon.GenerateRandomBytes(),
		Sender:   testscommon.GenerateRandomBytes()}
	ret := ap.ProcessAccounts([]*schema.Transaction{tx}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 0)
}

func TestAccountsProcessor_ProcessAccounts_FourAddresses_TwoIdentical_ExpectTwoAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountStub{}},
		&mock.PubKeyConverterStub{})

	tx1 := &schema.Transaction{
		Receiver: []byte("adr1"),
		Sender:   []byte("adr2")}
	tx2 := &schema.Transaction{
		Receiver: []byte("adr2"),
		Sender:   []byte("adr1")}

	ret := ap.ProcessAccounts([]*schema.Transaction{tx1, tx2}, []*schema.SCResult{}, []*schema.Receipt{})

	require.Len(t, ret, 2)
	for idx, account := range ret {
		require.Equal(t, account.Address, []byte("erd1addr"+strconv.Itoa(idx)))
		require.Equal(t, account.Balance, big.NewInt(int64(idx+1)).Bytes())
		require.Equal(t, account.Nonce, int64(idx+1))
	}
}

func TestAccountsProcessor_ProcessAccounts_SevenAddresses_ExpectSevenAccounts(t *testing.T) {
	ap, _ := accounts.NewAccountsProcessor(
		&mock.ShardCoordinatorMock{},
		&mock.AccountsAdapterStub{UserAccountHandler: &mock.UserAccountStub{}},
		&mock.PubKeyConverterStub{})

	tx1 := &schema.Transaction{
		Receiver: []byte("r1tx1"),
		Sender:   []byte("s1tx1"),
	}
	tx2 := &schema.Transaction{
		Receiver: []byte("r1tx2"),
		Sender:   []byte("s1tx2"),
	}
	txs := []*schema.Transaction{tx1, tx2}

	scr := &schema.SCResult{
		Receiver: []byte("r1scr1"),
		Sender:   []byte("s1scr1"),
	}
	scrs := []*schema.SCResult{scr}

	receipt := &schema.Receipt{
		Sender: []byte("s1receipt1"),
	}
	receipts := []*schema.Receipt{receipt}

	ret := ap.ProcessAccounts(txs, scrs, receipts)

	require.Len(t, ret, 7)
	for idx, account := range ret {
		require.Equal(t, account.Address, []byte("erd1addr"+strconv.Itoa(idx)))
		require.Equal(t, account.Balance, big.NewInt(int64(idx+1)).Bytes())
		require.Equal(t, account.Nonce, int64(idx+1))
	}
}
