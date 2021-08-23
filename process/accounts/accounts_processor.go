package accounts

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data"
)

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
	processedReceipts []*schema.Receipt) ([]*schema.AccountBalanceUpdate, error) {

	addresses := ap.getAllAddresses(processedTxs, processedSCRs, processedReceipts)
	_ = addresses

	sender := processedTxs[1].Sender
	_ = processedSCRs[1].Receiver
	_ = processedReceipts[1].Sender

	if ap.shardCoordinator.SelfShardID() == ap.shardCoordinator.ComputeShardID(sender) {
		// Add to map
	}
	// SAME for receiver

	accountSender, _ := ap.accounts.LoadAccount(sender)

	newAcc := accountSender.(data.UserAccountHandler)

	newAcc.GetBalance()

	var myMap map[string]struct{}
	myMap["da"] = struct{}{}

	return nil, nil
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
	if ap.shardCoordinator.SelfShardID() == ap.shardCoordinator.ComputeShardID(address) {
		addresses[string(address)] = struct{}{}
	}
}
