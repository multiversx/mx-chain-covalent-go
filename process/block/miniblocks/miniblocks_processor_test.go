package miniblocks_test

import (
	"errors"
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/block/miniblocks"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMiniBlocksProcessor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        func() (hashing.Hasher, marshal.Marshalizer)
		expectedErr error
	}{
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return nil, &mock.MarshallerStub{}
			},
			expectedErr: covalent.ErrNilHasher,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return &mock.HasherStub{}, nil
			},
			expectedErr: covalent.ErrNilMarshaller,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return &mock.HasherStub{}, &mock.MarshallerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := miniblocks.NewMiniBlocksProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}

func TestMiniBlocksProcessor_ProcessMiniBlocks(t *testing.T) {
	mbp, _ := miniblocks.NewMiniBlocksProcessor(&mock.HasherStub{}, &mock.MarshallerStub{})

	header := &block.Header{TimeStamp: 123}
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			TxHashes:        [][]byte{[]byte("x"), []byte("y")},
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            3},
		{
			TxHashes:        [][]byte{[]byte("y"), []byte("z")},
			ReceiverShardID: 4,
			SenderShardID:   5,
			Type:            6},
	}}

	ret, _ := mbp.ProcessMiniBlocks(header, body)

	require.Len(t, ret, 2)

	require.Equal(t, ret[0].Hash, []byte("ok"))
	require.Equal(t, ret[0].TxHashes, [][]byte{[]byte("x"), []byte("y")})
	require.Equal(t, ret[0].Timestamp, int64(123))
	require.Equal(t, ret[0].ReceiverShardID, int32(1))
	require.Equal(t, ret[0].SenderShardID, int32(2))
	require.Equal(t, ret[0].Type, int32(3))

	require.Equal(t, ret[1].Hash, []byte("ok"))
	require.Equal(t, ret[1].TxHashes, [][]byte{[]byte("y"), []byte("z")})
	require.Equal(t, ret[1].Timestamp, int64(123))
	require.Equal(t, ret[1].ReceiverShardID, int32(4))
	require.Equal(t, ret[1].SenderShardID, int32(5))
	require.Equal(t, ret[1].Type, int32(6))
}

func TestMiniBlocksProcessor_ProcessMiniBlocks_InvalidMarshaller_ExpectZeroMBProcessed(t *testing.T) {
	mbp, _ := miniblocks.NewMiniBlocksProcessor(
		&mock.HasherStub{},
		&mock.MarshallerStub{
			MarshalCalled: func(obj interface{}) ([]byte, error) {
				return nil, errors.New("error marshaller stub")
			},
		})

	header := &block.Header{TimeStamp: 123}
	body := &block.Body{MiniBlocks: []*block.MiniBlock{
		{
			ReceiverShardID: 1,
			SenderShardID:   2,
			Type:            3},
		{
			ReceiverShardID: 4,
			SenderShardID:   5,
			Type:            6},
	}}

	ret, err := mbp.ProcessMiniBlocks(header, body)

	require.Equal(t, err, nil)
	require.Len(t, ret, 0)
}
