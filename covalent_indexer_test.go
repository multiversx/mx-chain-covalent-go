package covalent_test

import (
	"net/http"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go"
	"github.com/ElrondNetwork/covalent-indexer-go/process/ws"
	"github.com/ElrondNetwork/covalent-indexer-go/schema"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon"
	"github.com/ElrondNetwork/covalent-indexer-go/testscommon/mock"
	"github.com/ElrondNetwork/elrond-vm-common/atomic"
	"github.com/stretchr/testify/require"
)

func TestCovalentIndexer_sendBlockResultToCovalent_InvalidBlock_ExpectNoSentMessage(t *testing.T) {
	ci := covalent.NewCovalentDataIndexer(nil, &http.Server{
		Addr: "http://localhost:8080",
	})
	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}
	ci.SetWSSender(wss)

	ci.SendBlockResultToCovalent(nil)
	require.False(t, called.IsSet())
}

func TestCovalentIndexer_sendBlockResultToCovalent_ExpectSentMessage(t *testing.T) {
	ci := covalent.NewCovalentDataIndexer(nil, &http.Server{
		Addr: "http://localhost:8080",
	})
	called := atomic.Flag{}
	called.Unset()
	wss := &ws.WSSender{
		Conn: &mock.WSConnStub{
			WriteMessageCalled: func(messageType int, data []byte) error {
				called.Set()
				return nil
			},
		},
	}

	ci.SetWSSender(wss)

	blockRes := generateRandomValidBlockResult()
	ci.SendBlockResultToCovalent(blockRes)
	require.True(t, called.IsSet())
}

func generateRandomValidBlockResult() *schema.BlockResult {
	block := &schema.Block{
		Hash:          testscommon.GenerateRandomFixedBytes(32),
		StateRootHash: testscommon.GenerateRandomFixedBytes(32),
	}

	return &schema.BlockResult{
		Block: block,
	}
}
