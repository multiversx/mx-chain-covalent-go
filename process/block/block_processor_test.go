package block

import (
	"encoding/json"
	"errors"
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/elrond-go-core/data"
	erdBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBlockProcessor_NewBlockProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (marshal.Marshalizer, process.MiniBlockHandler)
		expectedErr error
	}{
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return nil, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: covalent.ErrNilMarshalizer,
		},
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.MarshallerStub{}, nil
			},
			expectedErr: covalent.ErrNilMiniBlockHandler,
		},
		{
			args: func() (marshal.Marshalizer, process.MiniBlockHandler) {
				return &mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := NewBlockProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_ProcessBlock_InvalidBodyAndHeaderMarshaller_ExpectProcessError(t *testing.T) {
	errMarshallHeader := errors.New("err header marshall")
	errMarshallBody := errors.New("err body marshall")

	tests := []struct {
		Marshaller  func(obj interface{}) ([]byte, error)
		expectedErr error
	}{
		{
			Marshaller: func(obj interface{}) ([]byte, error) {
				_, ok := obj.(*erdBlock.Header)
				if ok {
					return nil, errMarshallHeader
				}
				return json.Marshal(obj)
			},
			expectedErr: errMarshallHeader,
		},
		{
			Marshaller: func(obj interface{}) ([]byte, error) {
				_, ok := obj.(*erdBlock.Body)
				if ok {
					return nil, errMarshallBody
				}
				return json.Marshal(obj)
			},
			expectedErr: errMarshallBody,
		},
	}

	for _, currTest := range tests {
		bp, _ := NewBlockProcessor(
			&mock.MarshallerStub{
				MarshalCalled: currTest.Marshaller},
			&mock.MiniBlockHandlerStub{})

		args := getInitializedArgs(false)
		_, err := bp.ProcessBlock(args)

		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestBlockProcessor_ProcessBlock_InvalidBody_ExpectErrBlockBodyAssertion(t *testing.T) {
	bp, _ := NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	args := getInitializedArgs(false)
	args.Body = nil
	_, err := bp.ProcessBlock(args)

	require.Equal(t, err, covalent.ErrBlockBodyAssertion)
}

func TestNewBlockProcessor_ProcessBlock_InvalidMBHandler_ExpectErr(t *testing.T) {
	errMBHandler := errors.New("error mb handler")

	bp, _ := NewBlockProcessor(
		&mock.MarshallerStub{},
		&mock.MiniBlockHandlerStub{
			ProcessMiniBlockCalled: func(headerHash []byte, header data.HeaderHandler, body *erdBlock.Body) ([]*schema.MiniBlock, error) {
				return nil, errMBHandler
			}})

	args := getInitializedArgs(false)
	_, err := bp.ProcessBlock(args)

	require.Equal(t, err, errMBHandler)
}

func TestNewBlockProcessor_ProcessBlock_NoSigners_ExpectDefaultProposerIndex(t *testing.T) {
	bp, _ := NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	args := getInitializedArgs(false)
	args.SignersIndexes = nil
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Proposer, DefaultProposerIndex)
}

func TestBlockProcessor_ProcessBlock(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(false)
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Nonce, int64(args.Header.GetNonce()))
	require.Equal(t, ret.Round, int64(args.Header.GetRound()))
	require.Equal(t, ret.Epoch, int32(args.Header.GetEpoch()))
	require.Equal(t, ret.Hash, args.HeaderHash)
	require.Equal(t, ret.NotarizedBlocksHashes, utility.StrSliceToBytesSlice(args.NotarizedHeadersHashes)) //TODO: test this utlity
	require.Equal(t, ret.Proposer, int64(args.SignersIndexes[0]))
	require.Equal(t, ret.Validators, utility.UIntSliceToIntSlice(args.SignersIndexes))
	require.Equal(t, ret.PubKeysBitmap, args.Header.GetPubKeysBitmap())
	require.Equal(t, ret.Size, int64(485))
	require.Equal(t, ret.Timestamp, int64(args.Header.GetTimeStamp()))
	require.Equal(t, ret.StateRootHash, args.Header.GetRootHash())
	require.Equal(t, ret.PrevHash, args.Header.GetPrevHash())
	require.Equal(t, ret.ShardID, int32(args.Header.GetShardID()))
	require.Equal(t, ret.TxCount, int32(args.Header.GetTxCount()))
	require.Equal(t, ret.AccumulatedFees, args.Header.GetAccumulatedFees().Bytes())
	require.Equal(t, ret.DeveloperFees, args.Header.GetDeveloperFees().Bytes())

	require.Equal(t, ret.EpochStartInfo, (*schema.EpochStartInfo)(nil))
}

func TestBlockProcessor_ProcessMetaBlock(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})
	args := getInitializedArgs(true)
	ret, _ := bp.ProcessBlock(args)

	require.Equal(t, ret.Nonce, int64(args.Header.GetNonce()))
	require.Equal(t, ret.Round, int64(args.Header.GetRound()))
	require.Equal(t, ret.Epoch, int32(args.Header.GetEpoch()))
	require.Equal(t, ret.Hash, args.HeaderHash)
	require.Equal(t, ret.NotarizedBlocksHashes, utility.StrSliceToBytesSlice(args.NotarizedHeadersHashes)) //TODO: test this utlity
	require.Equal(t, ret.Proposer, int64(args.SignersIndexes[0]))
	require.Equal(t, ret.Validators, utility.UIntSliceToIntSlice(args.SignersIndexes))
	require.Equal(t, ret.PubKeysBitmap, args.Header.GetPubKeysBitmap())
	require.Equal(t, ret.Timestamp, int64(args.Header.GetTimeStamp()))
	require.Equal(t, ret.StateRootHash, args.Header.GetRootHash())
	require.Equal(t, ret.PrevHash, args.Header.GetPrevHash())
	require.Equal(t, ret.ShardID, int32(args.Header.GetShardID()))
	require.Equal(t, ret.TxCount, int32(args.Header.GetTxCount()))
	require.Equal(t, ret.AccumulatedFees, args.Header.GetAccumulatedFees().Bytes())
	require.Equal(t, ret.DeveloperFees, args.Header.GetDeveloperFees().Bytes())

	metaBlockEconomics := args.Header.(*erdBlock.MetaBlock).GetEpochStart().Economics

	require.Equal(t, ret.EpochStartInfo.TotalSupply, metaBlockEconomics.TotalSupply.Bytes())
	require.Equal(t, ret.EpochStartInfo.TotalToDistribute, metaBlockEconomics.TotalToDistribute.Bytes())
	require.Equal(t, ret.EpochStartInfo.TotalNewlyMinted, metaBlockEconomics.TotalNewlyMinted.Bytes())
	require.Equal(t, ret.EpochStartInfo.RewardsPerBlock, metaBlockEconomics.RewardsPerBlock.Bytes())
	require.Equal(t, ret.EpochStartInfo.RewardsForProtocolSustainability, metaBlockEconomics.RewardsForProtocolSustainability.Bytes())
	require.Equal(t, ret.EpochStartInfo.NodePrice, metaBlockEconomics.NodePrice.Bytes())
	require.Equal(t, ret.EpochStartInfo.PrevEpochStartRound, int32(metaBlockEconomics.PrevEpochStartRound))
	require.Equal(t, ret.EpochStartInfo.PrevEpochStartHash, metaBlockEconomics.PrevEpochStartHash)
}

func TestBlockProcessor_ProcessMetaBlock_NotStartOfEpochBlock_ExpectNilEpochStartInfo(t *testing.T) {
	bp, _ := NewBlockProcessor(&mock.MarshallerStub{}, &mock.MiniBlockHandlerStub{})

	metaBlockHeader := getInitializedMetaBlockHeader()
	metaBlockHeader.EpochStart.LastFinalizedHeaders = nil

	ret, _ := bp.ProcessBlock(&indexer.ArgsSaveBlockData{
		Header: metaBlockHeader,
		Body:   &erdBlock.Body{}})

	require.Equal(t, ret.EpochStartInfo, (*schema.EpochStartInfo)(nil))
}

func getInitializedArgs(metaBlock bool) *indexer.ArgsSaveBlockData {
	var header data.HeaderHandler

	if metaBlock {
		header = getInitializedMetaBlockHeader()
	} else {
		header = getInitialisedHeader()
	}

	return &indexer.ArgsSaveBlockData{
		HeaderHash:             []byte("header hash"),
		Body:                   &erdBlock.Body{},
		Header:                 header,
		SignersIndexes:         []uint64{1, 2, 3},
		NotarizedHeadersHashes: []string{"h1", "h2"},
		TransactionsPool:       nil,
	}
}

func getInitializedMetaBlockHeader() *erdBlock.MetaBlock {
	return &erdBlock.MetaBlock{
		Nonce:           1,
		Epoch:           2,
		Round:           3,
		TimeStamp:       4,
		LeaderSignature: []byte("meta leader signature"),
		PubKeysBitmap:   []byte("meta pub keys bitmap"),
		PrevHash:        []byte("meta prev hash"),
		PrevRandSeed:    []byte("meta prev rand seed"),
		RandSeed:        []byte("meta rand seed"),
		RootHash:        []byte("meta root hash"),
		EpochStart: erdBlock.EpochStart{
			LastFinalizedHeaders: []erdBlock.EpochStartShardData{{}},
			Economics: erdBlock.Economics{
				TotalSupply:                      big.NewInt(5),
				TotalToDistribute:                big.NewInt(6),
				TotalNewlyMinted:                 big.NewInt(7),
				RewardsPerBlock:                  big.NewInt(8),
				RewardsForProtocolSustainability: big.NewInt(9),
				NodePrice:                        big.NewInt(10),
				PrevEpochStartRound:              11,
				PrevEpochStartHash:               []byte("meta prev epoch hash"),
			},
		},
		ChainID:         []byte("meta chain id"),
		AccumulatedFees: big.NewInt(11),
		DeveloperFees:   big.NewInt(12),
		TxCount:         13,
	}
}

func getInitialisedHeader() *erdBlock.Header {
	return &erdBlock.Header{
		Nonce:              1,
		PrevHash:           []byte("prev hash"),
		PrevRandSeed:       []byte("prev rand seed"),
		RandSeed:           []byte("rand seed"),
		PubKeysBitmap:      []byte("pub keys bitmap"),
		ShardID:            2,
		TimeStamp:          3,
		Round:              4,
		Epoch:              5,
		BlockBodyType:      6,
		LeaderSignature:    []byte("leader signature"),
		RootHash:           []byte("root hash"),
		TxCount:            7,
		EpochStartMetaHash: []byte("epoch start meta hash"),
		ReceiptsHash:       []byte("receipts hash"),
		ChainID:            []byte("chain id"),
		AccumulatedFees:    big.NewInt(8),
		DeveloperFees:      big.NewInt(9),
	}
}
