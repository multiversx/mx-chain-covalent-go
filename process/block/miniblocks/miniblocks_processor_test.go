package miniblocks

import (
	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/mock"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockProcessor_NewBlockProcessor(t *testing.T) {
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
			expectedErr: covalent.ErrNilMarshalizer,
		},
		{
			args: func() (hashing.Hasher, marshal.Marshalizer) {
				return &mock.HasherStub{}, &mock.MarshallerStub{}
			},
			expectedErr: nil,
		},
	}

	for _, currTest := range tests {
		_, err := NewMiniBlocksProcessor(currTest.args())
		require.Equal(t, currTest.expectedErr, err)
	}
}
