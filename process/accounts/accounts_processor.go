package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("process/accounts")

type accountsProcessor struct {
	shardCoordinator process.ShardCoordinator
	pubKeyConverter  core.PubkeyConverter
	accounts         covalent.AccountsAdapter
}

// NewAccountsProcessor creates a new instance of accounts processor
func NewAccountsProcessor(
	shardCoordinator process.ShardCoordinator,
	accounts covalent.AccountsAdapter,
	pubKeyConverter core.PubkeyConverter) (*accountsProcessor, error) {

	if check.IfNil(shardCoordinator) {
		return nil, covalent.ErrNilShardCoordinator
	}
	if check.IfNil(accounts) {
		return nil, covalent.ErrNilAccountsAdapter
	}
	if check.IfNil(pubKeyConverter) {
		return nil, covalent.ErrNilPubKeyConverter
	}

	return &accountsProcessor{
		accounts:         accounts,
		pubKeyConverter:  pubKeyConverter,
		shardCoordinator: shardCoordinator,
	}, nil
}

// ProcessAccounts converts accounts data to a specific structure defined by avro schema
func (ap *accountsProcessor) ProcessAccounts(
	processedTxs []*schema.Transaction,
	processedSCRs []*schema.SCResult,
	processedReceipts []*schema.Receipt) []*schema.AccountBalanceUpdate {

	addresses := ap.getAllAddresses(processedTxs, processedSCRs, processedReceipts)
	accounts := make([]*schema.AccountBalanceUpdate, 0)

	for address, _ := range addresses {
		account, err := ap.processAccount(address)
		if err != nil {
			log.Warn("cannot get address account", "address", address, "error", err)
			continue
		}

		accounts = append(accounts, account)
	}

	return accounts
}

func (ap *accountsProcessor) getAllAddresses(
	processedTxs []*schema.Transaction,
	processedSCRs []*schema.SCResult,
	processedReceipts []*schema.Receipt) map[string]struct{} {
	addresses := make(map[string]struct{})

	for _, tx := range processedTxs {
		ap.addAddressIfInSelfShard(addresses, tx.Sender)
		ap.addAddressIfInSelfShard(addresses, tx.Receiver)
	}

	for _, scr := range processedSCRs {
		ap.addAddressIfInSelfShard(addresses, scr.Sender)
		ap.addAddressIfInSelfShard(addresses, scr.Receiver)
	}

	for _, receipt := range processedReceipts {
		ap.addAddressIfInSelfShard(addresses, receipt.Sender)
	}

	return addresses
}

func (ap *accountsProcessor) addAddressIfInSelfShard(addresses map[string]struct{}, address []byte) {
	if ap.shardCoordinator.SelfId() == ap.shardCoordinator.ComputeId(address) {
		addresses[string(address)] = struct{}{}
	}
}

func (ap *accountsProcessor) processAccount(address string) (*schema.AccountBalanceUpdate, error) {
	acc, err := ap.accounts.LoadAccount([]byte(address))
	if err != nil || check.IfNil(acc) {
		return nil, err
	}

	account, castOk := acc.(data.UserAccountHandler)
	if !castOk {
		return nil, covalent.ErrCannotCastAccountHandlerToUserAccount
	}

	return &schema.AccountBalanceUpdate{
		Address: utility.EncodePubKey(ap.pubKeyConverter, account.AddressBytes()),
		Balance: utility.GetBytes(account.GetBalance()),
		Nonce:   int64(account.GetNonce()),
	}, nil
}
