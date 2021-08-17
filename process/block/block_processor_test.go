package block

import (
	"errors"
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBlockProcessor_NewBlockProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler)
		expectedErr error
	}{
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.HasherMock{}, nil, nil
			},
			expectedErr: covalent.ErrNilMarshalizer,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return nil, &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: covalent.ErrNilHasher,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.HasherMock{}, &mock.MarshallerStub{}, nil
			},
			expectedErr: covalent.ErrNilMiniBlockHandler,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.HasherMock{}, &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := NewBlockProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_Marshal(t *testing.T) {

	expectedError := errors.New("expectedError")
	marshaller := &mock.MarshallerStub{
		MarshalCalled: func(obj interface{}) ([]byte, error) {
			return nil, expectedError
		},
	}
	_ = marshaller

}

func TestBlockProcessor_ProcessBlock(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.HasherMock{}, &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(false)
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Nonce, int64(nonce))
	require.Equal(t, ret.Round, int64(round))
	require.Equal(t, ret.Epoch, int32(epoch))
	require.Equal(t, ret.Hash, headerHash)
	require.Equal(t, ret.NotarizedBlocksHashes, utility.StrSliceToBytesSlice(notarizedHeadersHashes)) //TODO: test this utlity
	require.Equal(t, ret.Proposer, int64(signersIndexes[0]))
	require.Equal(t, ret.Validators, utility.UIntSliceToIntSlice(signersIndexes))
	require.Equal(t, ret.PubKeysBitmap, pubKeysBitmap)
	require.Equal(t, ret.Size, int64(473))
	require.Equal(t, ret.SizeTxs, int64(0))
	require.Equal(t, ret.Timestamp, int64(timeStamp))
	require.Equal(t, ret.StateRootHash, rootHash)
	require.Equal(t, ret.PrevHash, prevHash)
	require.Equal(t, ret.ShardID, int32(shardID))
	require.Equal(t, ret.TxCount, int32(txCount))
	require.Equal(t, ret.AccumulatedFees, accumulatedFees.Bytes())
	require.Equal(t, ret.DeveloperFees, developerFees.Bytes())

	require.Equal(t, ret.EpochStartInfo, (*schema.EpochStartInfo)(nil))
}

func TestBlockProcessor_ProcessMetaBlock(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.HasherMock{}, &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(true)
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Nonce, int64(nonce))
	require.Equal(t, ret.Round, int64(round))
	require.Equal(t, ret.Epoch, int32(epoch))
	require.Equal(t, ret.Hash, headerHash)
	require.Equal(t, ret.NotarizedBlocksHashes, utility.StrSliceToBytesSlice(notarizedHeadersHashes)) //TODO: test this utlity
	require.Equal(t, ret.Proposer, int64(signersIndexes[0]))
	require.Equal(t, ret.Validators, utility.UIntSliceToIntSlice(signersIndexes))
	require.Equal(t, ret.PubKeysBitmap, pubKeysBitmap)
	require.Equal(t, ret.SizeTxs, int64(0))
	require.Equal(t, ret.Timestamp, int64(timeStamp))
	require.Equal(t, ret.StateRootHash, rootHash)
	require.Equal(t, ret.PrevHash, prevHash)
	require.Equal(t, uint32(ret.ShardID), core.MetachainShardId)
	require.Equal(t, ret.TxCount, int32(txCount))
	require.Equal(t, ret.AccumulatedFees, accumulatedFees.Bytes())
	require.Equal(t, ret.DeveloperFees, developerFees.Bytes())

	require.Equal(t, ret.EpochStartInfo.TotalSupply, economicsTotalSupply.Bytes())
	require.Equal(t, ret.EpochStartInfo.TotalToDistribute, economicsTotalToDistribute.Bytes())
	require.Equal(t, ret.EpochStartInfo.TotalNewlyMinted, economicsTotalNewlyMinted.Bytes())
	require.Equal(t, ret.EpochStartInfo.RewardsPerBlock, economicsRewardsPerBlock.Bytes())
	require.Equal(t, ret.EpochStartInfo.RewardsForProtocolSustainability, economicsRewardsForProtocolSustainability.Bytes())
	require.Equal(t, ret.EpochStartInfo.NodePrice, economicsNodePrice.Bytes())
	require.Equal(t, ret.EpochStartInfo.PrevEpochStartRound, int64(economicsPrevEpochStartRound))
	require.Equal(t, ret.EpochStartInfo.PrevEpochStartHash, economicsPrevEpochStartHash)
}

func TestBlockProcessor_ProcessMetaBlock_NotStartOfEpochBlock_ExpectNilEpochStartInfo(t *testing.T) {
	bp, _ := NewBlockProcessor(&mock.HasherMock{}, &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	metaBlockHeader := getInitializedMetaBlockHeader()
	metaBlockHeader.EpochStart.LastFinalizedHeaders = nil

	ret, _ := bp.ProcessBlock(&indexer.ArgsSaveBlockData{
		Header: metaBlockHeader,
		Body:   &erdBlock.Body{}})

	require.Equal(t, ret.EpochStartInfo, (*schema.EpochStartInfo)(nil))
}

var headerHash = []byte("header hash")
var signersIndexes = []uint64{444, 333, 222}
var notarizedHeadersHashes = []string{"h1", "h2"}

var nonce = uint64(1)
var prevHash = []byte("prev hash")
var prevRandSeed = []byte("prev rand seed")
var randSeed = []byte("rand seed")
var pubKeysBitmap = []byte("pub keys bitmap")
var shardID = uint32(2)
var timeStamp = uint64(3)
var round = uint64(4)
var epoch = uint32(5)
var blockBodyType = erdBlock.Type(6)
var leaderSignature = []byte("rand seed")
var rootHash = []byte("root hash")
var txCount = uint32(7)
var epochStartMetaHash = []byte("epoch start meta hash")
var receiptsHash = []byte("receipts hash")
var chainID = []byte("chain id")
var accumulatedFees = big.NewInt(8)
var developerFees = big.NewInt(9)

var economicsTotalSupply = big.NewInt(10)
var economicsTotalToDistribute = big.NewInt(11)
var economicsTotalNewlyMinted = big.NewInt(12)
var economicsRewardsPerBlock = big.NewInt(13)
var economicsRewardsForProtocolSustainability = big.NewInt(14)
var economicsNodePrice = big.NewInt(15)
var economicsPrevEpochStartRound = uint64(123)
var economicsPrevEpochStartHash = []byte("econ prev epoch start hash")

func getInitializedArgs(metaBlock bool) *indexer.ArgsSaveBlockData {
	var header data.HeaderHandler

	if metaBlock {
		header = getInitializedMetaBlockHeader()
	} else {
		header = getInitialisedHeader()
	}

	return &indexer.ArgsSaveBlockData{
		HeaderHash:             headerHash,
		Body:                   &erdBlock.Body{},
		Header:                 header,
		SignersIndexes:         signersIndexes,
		NotarizedHeadersHashes: notarizedHeadersHashes,
		TransactionsPool:       nil,
	}
}

func getInitializedMetaBlockHeader() *erdBlock.MetaBlock {
	return &erdBlock.MetaBlock{
		Nonce:                  nonce,
		Epoch:                  epoch,
		Round:                  round,
		TimeStamp:              timeStamp,
		ShardInfo:              nil,
		PeerInfo:               nil,
		Signature:              nil,
		LeaderSignature:        leaderSignature,
		PubKeysBitmap:          pubKeysBitmap,
		PrevHash:               prevHash,
		PrevRandSeed:           prevRandSeed,
		RandSeed:               randSeed,
		RootHash:               rootHash,
		ValidatorStatsRootHash: nil,
		MiniBlockHeaders:       nil,
		ReceiptsHash:           nil,
		EpochStart: erdBlock.EpochStart{
			LastFinalizedHeaders: []erdBlock.EpochStartShardData{{}},
			Economics: erdBlock.Economics{
				TotalSupply:                      economicsTotalSupply,
				TotalToDistribute:                economicsTotalToDistribute,
				TotalNewlyMinted:                 economicsTotalNewlyMinted,
				RewardsPerBlock:                  economicsRewardsPerBlock,
				RewardsForProtocolSustainability: economicsRewardsForProtocolSustainability,
				NodePrice:                        economicsNodePrice,
				PrevEpochStartRound:              economicsPrevEpochStartRound,
				PrevEpochStartHash:               economicsPrevEpochStartHash,
			},
		},
		ChainID:                nil,
		SoftwareVersion:        nil,
		AccumulatedFees:        accumulatedFees,
		AccumulatedFeesInEpoch: nil,
		DeveloperFees:          developerFees,
		DevFeesInEpoch:         nil,
		TxCount:                txCount,
		Reserved:               nil,
	}
}

func getInitialisedHeader() *erdBlock.Header {
	return &erdBlock.Header{
		Nonce:              nonce,
		PrevHash:           prevHash,
		PrevRandSeed:       prevRandSeed,
		RandSeed:           randSeed,
		PubKeysBitmap:      pubKeysBitmap,
		ShardID:            shardID,
		TimeStamp:          timeStamp,
		Round:              round,
		Epoch:              epoch,
		BlockBodyType:      blockBodyType,
		Signature:          nil,
		LeaderSignature:    leaderSignature,
		MiniBlockHeaders:   nil,
		PeerChanges:        nil,
		RootHash:           rootHash,
		MetaBlockHashes:    nil,
		TxCount:            txCount,
		EpochStartMetaHash: epochStartMetaHash,
		ReceiptsHash:       receiptsHash,
		ChainID:            chainID,
		AccumulatedFees:    accumulatedFees,
		DeveloperFees:      developerFees,
	}
}
