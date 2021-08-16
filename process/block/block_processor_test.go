package block

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/covalent-indexer-go/process"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBlockProcessor(t *testing.T) {
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
