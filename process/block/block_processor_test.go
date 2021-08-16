package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
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
				return nil, &mock.MarshalizerMock{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: covalent.ErrNilHasher,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.HasherMock{}, &mock.MarshalizerMock{}, nil
			},
			expectedErr: covalent.ErrNilMiniBlockHandler,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.HasherMock{}, &mock.MarshalizerMock{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := NewBlockProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_ProcessBlock(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.HasherMock{}, &mock.MarshalizerMock{}, &mock.MiniBlockHandlerStub{})

	args := getInitializedArgs()

	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Nonce, int64(nonce))
	require.Equal(t, ret.Round, int64(round))
	require.Equal(t, ret.Epoch, int32(epoch))
	require.Equal(t, ret.Hash, headerHash)
	require.Equal(t, ret.NotarizedBlocksHashes, utility.StrSliceToBytesSlice(notarizedHeadersHashes)) //TODO: test this utlity
	require.Equal(t, ret.Proposer, int64(signersIndexes[0]))
	require.Equal(t, ret.Validators, utility.UIntSliceToIntSlice(signersIndexes))
	require.Equal(t, ret.PubKeysBitmap, pubKeysBitmap)
	//require.Equal(t, ret.Size, int64(5))
	require.Equal(t, ret.SizeTxs, int64(0))
	require.Equal(t, ret.Timestamp, int64(timeStamp))
	require.Equal(t, ret.StateRootHash, rootHash)
	require.Equal(t, ret.PrevHash, prevHash)
	require.Equal(t, ret.ShardID, int32(shardID))
	require.Equal(t, ret.TxCount, int32(txCount))
	require.Equal(t, ret.AccumulatedFees, accumulatedFees.Bytes())
	require.Equal(t, ret.DeveloperFees, developerFees.Bytes())
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

func getInitializedArgs() *indexer.ArgsSaveBlockData {
	return &indexer.ArgsSaveBlockData{
		HeaderHash:             headerHash,
		Body:                   &erdBlock.Body{},
		Header:                 getInitialisedHeader(),
		SignersIndexes:         signersIndexes,
		NotarizedHeadersHashes: notarizedHeadersHashes,
		TransactionsPool:       nil,
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
