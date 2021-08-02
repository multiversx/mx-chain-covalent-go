package schema

import "github.com/elodina/go-avro"

type BlockResult struct {
	Block        *Block
	Transactions []*Transaction
	SCResults    []*SCResult
	Receipts     []*Receipt
	Logs         []*Log
	StateChanges []*AccountBalanceUpdate
}

func NewBlockResult() *BlockResult {
	return &BlockResult{
		Block:        NewBlock(),
		Transactions: make([]*Transaction, 0),
		SCResults:    make([]*SCResult, 0),
		Receipts:     make([]*Receipt, 0),
		Logs:         make([]*Log, 0),
		StateChanges: make([]*AccountBalanceUpdate, 0),
	}
}

func (o *BlockResult) Schema() avro.Schema {
	if _BlockResult_schema_err != nil {
		panic(_BlockResult_schema_err)
	}
	return _BlockResult_schema
}

type Block struct {
	Nonce                 int64
	Round                 int64
	Epoch                 int32
	Hash                  string
	MiniBlocks            []*MiniBlock
	NotarizedBlocksHashes []string
	Proposer              int64
	Validators            []int64
	PubKeysBitmap         []byte
	Size                  int64
	SizeTxs               int64
	Timestamp             int64
	StateRootHash         string
	PrevHash              string
	ShardID               int32
	TxCount               int32
	AccumulatedFees       string
	DeveloperFees         string
	EpochStartBlock       bool
	EpochStartInfo        *EpochStartInfo
}

func NewBlock() *Block {
	return &Block{
		MiniBlocks:            make([]*MiniBlock, 0),
		NotarizedBlocksHashes: make([]string, 0),
		Validators:            make([]int64, 0),
		PubKeysBitmap:         []byte{},
		EpochStartInfo:        NewEpochStartInfo(),
	}
}

func (o *Block) Schema() avro.Schema {
	if _Block_schema_err != nil {
		panic(_Block_schema_err)
	}
	return _Block_schema
}

type MiniBlock struct {
	Hash              string
	SenderShardID     int32
	ReceiverShardID   int32
	SenderBlockHash   string
	ReceiverBlockHash string
	Type              string
	Timestamp         int64
}

func NewMiniBlock() *MiniBlock {
	return &MiniBlock{}
}

func (o *MiniBlock) Schema() avro.Schema {
	if _MiniBlock_schema_err != nil {
		panic(_MiniBlock_schema_err)
	}
	return _MiniBlock_schema
}

type EpochStartInfo struct {
	TotalSupply                      string
	TotalToDistribute                string
	TotalNewlyMinted                 string
	RewardsPerBlock                  string
	RewardsForProtocolSustainability string
	NodePrice                        string
	PrevEpochStartRound              int32
	PrevEpochStartHash               string
}

func NewEpochStartInfo() *EpochStartInfo {
	return &EpochStartInfo{}
}

func (o *EpochStartInfo) Schema() avro.Schema {
	if _EpochStartInfo_schema_err != nil {
		panic(_EpochStartInfo_schema_err)
	}
	return _EpochStartInfo_schema
}

type Transaction struct {
	Hash             string
	MiniBlockHash    string
	BlockHash        string
	Nonce            int64
	Round            int64
	Value            string
	Receiver         string
	Sender           string
	ReceiverShard    int32
	SenderShard      int32
	GasPrice         int64
	GasLimit         int64
	Signature        string
	Timestamp        int64
	SenderUserName   []byte
	ReceiverUserName []byte
}

func NewTransaction() *Transaction {
	return &Transaction{
		SenderUserName:   []byte{},
		ReceiverUserName: []byte{},
	}
}

func (o *Transaction) Schema() avro.Schema {
	if _Transaction_schema_err != nil {
		panic(_Transaction_schema_err)
	}
	return _Transaction_schema
}

type SCResult struct {
	Hash           string
	Nonce          int64
	GasLimit       int64
	GasPrice       int64
	Value          string
	Sender         string
	Receiver       string
	RelayerAddr    string
	RelayedValue   string
	Code           string
	Data           []byte
	PrevTxHash     string
	OriginalTxHash string
	CallType       string
	CodeMetadata   []byte
	ReturnMessage  string
	Timestamp      int64
}

func NewSCResult() *SCResult {
	return &SCResult{
		Data:         []byte{},
		CodeMetadata: []byte{},
	}
}

func (o *SCResult) Schema() avro.Schema {
	if _SCResult_schema_err != nil {
		panic(_SCResult_schema_err)
	}
	return _SCResult_schema
}

type Receipt struct {
	Hash      string
	Value     string
	Sender    string
	Data      string
	TxHash    string
	Timestamp int64
}

func NewReceipt() *Receipt {
	return &Receipt{}
}

func (o *Receipt) Schema() avro.Schema {
	if _Receipt_schema_err != nil {
		panic(_Receipt_schema_err)
	}
	return _Receipt_schema
}

type Log struct {
	ID      string
	Address string
	Events  []*Event
}

func NewLog() *Log {
	return &Log{
		Events: make([]*Event, 0),
	}
}

func (o *Log) Schema() avro.Schema {
	if _Log_schema_err != nil {
		panic(_Log_schema_err)
	}
	return _Log_schema
}

type Event struct {
	Address    string
	Identifier string
	Topics     [][]byte
	Data       []byte
}

func NewEvent() *Event {
	return &Event{
		Topics: make([][]byte, 0),
		Data:   []byte{},
	}
}

func (o *Event) Schema() avro.Schema {
	if _Event_schema_err != nil {
		panic(_Event_schema_err)
	}
	return _Event_schema
}

type AccountBalanceUpdate struct {
	Address string
	Balance string
	Nonce   int64
}

func NewAccountBalanceUpdate() *AccountBalanceUpdate {
	return &AccountBalanceUpdate{}
}

func (o *AccountBalanceUpdate) Schema() avro.Schema {
	if _AccountBalanceUpdate_schema_err != nil {
		panic(_AccountBalanceUpdate_schema_err)
	}
	return _AccountBalanceUpdate_schema
}

// Generated by codegen. Please do not modify.
var _BlockResult_schema, _BlockResult_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "com.covalenthq.block.schema",
    "name": "BlockResult",
    "fields": [
        {
            "name": "Block",
            "type": {
                "type": "record",
                "name": "Block",
                "fields": [
                    {
                        "name": "Nonce",
                        "type": "long"
                    },
                    {
                        "name": "Round",
                        "type": "long"
                    },
                    {
                        "name": "Epoch",
                        "type": "int"
                    },
                    {
                        "name": "Hash",
                        "type": "string"
                    },
                    {
                        "name": "MiniBlocks",
                        "type": {
                            "type": "array",
                            "items": {
                                "type": "record",
                                "name": "MiniBlock",
                                "fields": [
                                    {
                                        "name": "Hash",
                                        "type": "string"
                                    },
                                    {
                                        "name": "SenderShardID",
                                        "type": "int"
                                    },
                                    {
                                        "name": "ReceiverShardID",
                                        "type": "int"
                                    },
                                    {
                                        "name": "SenderBlockHash",
                                        "type": "string"
                                    },
                                    {
                                        "name": "ReceiverBlockHash",
                                        "type": "string"
                                    },
                                    {
                                        "name": "Type",
                                        "type": "string"
                                    },
                                    {
                                        "name": "Timestamp",
                                        "type": "long"
                                    }
                                ]
                            }
                        }
                    },
                    {
                        "name": "NotarizedBlocksHashes",
                        "type": {
                            "type": "array",
                            "items": "string"
                        }
                    },
                    {
                        "name": "Proposer",
                        "type": "long"
                    },
                    {
                        "name": "Validators",
                        "type": {
                            "type": "array",
                            "items": "long"
                        }
                    },
                    {
                        "name": "PubKeysBitmap",
                        "type": "bytes"
                    },
                    {
                        "name": "Size",
                        "type": "long"
                    },
                    {
                        "name": "SizeTxs",
                        "type": "long"
                    },
                    {
                        "name": "Timestamp",
                        "type": "long"
                    },
                    {
                        "name": "StateRootHash",
                        "type": "string"
                    },
                    {
                        "name": "PrevHash",
                        "type": "string"
                    },
                    {
                        "name": "ShardID",
                        "type": "int"
                    },
                    {
                        "name": "TxCount",
                        "type": "int"
                    },
                    {
                        "name": "AccumulatedFees",
                        "type": "string"
                    },
                    {
                        "name": "DeveloperFees",
                        "type": "string"
                    },
                    {
                        "name": "EpochStartBlock",
                        "type": "boolean"
                    },
                    {
                        "name": "EpochStartInfo",
                        "type": {
                            "type": "record",
                            "name": "EpochStartInfo",
                            "fields": [
                                {
                                    "name": "TotalSupply",
                                    "type": "string"
                                },
                                {
                                    "name": "TotalToDistribute",
                                    "type": "string"
                                },
                                {
                                    "name": "TotalNewlyMinted",
                                    "type": "string"
                                },
                                {
                                    "name": "RewardsPerBlock",
                                    "type": "string"
                                },
                                {
                                    "name": "RewardsForProtocolSustainability",
                                    "type": "string"
                                },
                                {
                                    "name": "NodePrice",
                                    "type": "string"
                                },
                                {
                                    "name": "PrevEpochStartRound",
                                    "type": "int"
                                },
                                {
                                    "name": "PrevEpochStartHash",
                                    "type": "string"
                                }
                            ]
                        }
                    }
                ]
            }
        },
        {
            "name": "Transactions",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "Transaction",
                    "fields": [
                        {
                            "name": "Hash",
                            "type": "string"
                        },
                        {
                            "name": "MiniBlockHash",
                            "type": "string"
                        },
                        {
                            "name": "BlockHash",
                            "type": "string"
                        },
                        {
                            "name": "Nonce",
                            "type": "long"
                        },
                        {
                            "name": "Round",
                            "type": "long"
                        },
                        {
                            "name": "Value",
                            "type": "string"
                        },
                        {
                            "name": "Receiver",
                            "type": "string"
                        },
                        {
                            "name": "Sender",
                            "type": "string"
                        },
                        {
                            "name": "ReceiverShard",
                            "type": "int"
                        },
                        {
                            "name": "SenderShard",
                            "type": "int"
                        },
                        {
                            "name": "GasPrice",
                            "type": "long"
                        },
                        {
                            "name": "GasLimit",
                            "type": "long"
                        },
                        {
                            "name": "Signature",
                            "type": "string"
                        },
                        {
                            "name": "Timestamp",
                            "type": "long"
                        },
                        {
                            "name": "SenderUserName",
                            "type": "bytes"
                        },
                        {
                            "name": "ReceiverUserName",
                            "type": "bytes"
                        }
                    ]
                }
            }
        },
        {
            "name": "SCResults",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "SCResult",
                    "fields": [
                        {
                            "name": "Hash",
                            "type": "string"
                        },
                        {
                            "name": "Nonce",
                            "type": "long"
                        },
                        {
                            "name": "GasLimit",
                            "type": "long"
                        },
                        {
                            "name": "GasPrice",
                            "type": "long"
                        },
                        {
                            "name": "Value",
                            "type": "string"
                        },
                        {
                            "name": "Sender",
                            "type": "string"
                        },
                        {
                            "name": "Receiver",
                            "type": "string"
                        },
                        {
                            "name": "RelayerAddr",
                            "type": "string"
                        },
                        {
                            "name": "RelayedValue",
                            "type": "string"
                        },
                        {
                            "name": "Code",
                            "type": "string"
                        },
                        {
                            "name": "Data",
                            "type": "bytes"
                        },
                        {
                            "name": "PrevTxHash",
                            "type": "string"
                        },
                        {
                            "name": "OriginalTxHash",
                            "type": "string"
                        },
                        {
                            "name": "CallType",
                            "type": "string"
                        },
                        {
                            "name": "CodeMetadata",
                            "type": "bytes"
                        },
                        {
                            "name": "ReturnMessage",
                            "type": "string"
                        },
                        {
                            "name": "Timestamp",
                            "type": "long"
                        }
                    ]
                }
            }
        },
        {
            "name": "Receipts",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "Receipt",
                    "fields": [
                        {
                            "name": "Hash",
                            "type": "string"
                        },
                        {
                            "name": "Value",
                            "type": "string"
                        },
                        {
                            "name": "Sender",
                            "type": "string"
                        },
                        {
                            "name": "Data",
                            "type": "string"
                        },
                        {
                            "name": "TxHash",
                            "type": "string"
                        },
                        {
                            "name": "Timestamp",
                            "type": "long"
                        }
                    ]
                }
            }
        },
        {
            "name": "Logs",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "Log",
                    "fields": [
                        {
                            "name": "ID",
                            "type": "string"
                        },
                        {
                            "name": "Address",
                            "type": "string"
                        },
                        {
                            "name": "Events",
                            "type": {
                                "type": "array",
                                "items": {
                                    "type": "record",
                                    "name": "Event",
                                    "fields": [
                                        {
                                            "name": "Address",
                                            "type": "string"
                                        },
                                        {
                                            "name": "Identifier",
                                            "type": "string"
                                        },
                                        {
                                            "name": "Topics",
                                            "type": {
                                                "type": "array",
                                                "items": "bytes"
                                            }
                                        },
                                        {
                                            "name": "Data",
                                            "type": "bytes"
                                        }
                                    ]
                                }
                            }
                        }
                    ]
                }
            }
        },
        {
            "name": "StateChanges",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "AccountBalanceUpdate",
                    "fields": [
                        {
                            "name": "Address",
                            "type": "string"
                        },
                        {
                            "name": "Balance",
                            "type": "string"
                        },
                        {
                            "name": "Nonce",
                            "type": "long"
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Block_schema, _Block_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Block",
    "fields": [
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "Round",
            "type": "long"
        },
        {
            "name": "Epoch",
            "type": "int"
        },
        {
            "name": "Hash",
            "type": "string"
        },
        {
            "name": "MiniBlocks",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "MiniBlock",
                    "fields": [
                        {
                            "name": "Hash",
                            "type": "string"
                        },
                        {
                            "name": "SenderShardID",
                            "type": "int"
                        },
                        {
                            "name": "ReceiverShardID",
                            "type": "int"
                        },
                        {
                            "name": "SenderBlockHash",
                            "type": "string"
                        },
                        {
                            "name": "ReceiverBlockHash",
                            "type": "string"
                        },
                        {
                            "name": "Type",
                            "type": "string"
                        },
                        {
                            "name": "Timestamp",
                            "type": "long"
                        }
                    ]
                }
            }
        },
        {
            "name": "NotarizedBlocksHashes",
            "type": {
                "type": "array",
                "items": "string"
            }
        },
        {
            "name": "Proposer",
            "type": "long"
        },
        {
            "name": "Validators",
            "type": {
                "type": "array",
                "items": "long"
            }
        },
        {
            "name": "PubKeysBitmap",
            "type": "bytes"
        },
        {
            "name": "Size",
            "type": "long"
        },
        {
            "name": "SizeTxs",
            "type": "long"
        },
        {
            "name": "Timestamp",
            "type": "long"
        },
        {
            "name": "StateRootHash",
            "type": "string"
        },
        {
            "name": "PrevHash",
            "type": "string"
        },
        {
            "name": "ShardID",
            "type": "int"
        },
        {
            "name": "TxCount",
            "type": "int"
        },
        {
            "name": "AccumulatedFees",
            "type": "string"
        },
        {
            "name": "DeveloperFees",
            "type": "string"
        },
        {
            "name": "EpochStartBlock",
            "type": "boolean"
        },
        {
            "name": "EpochStartInfo",
            "type": {
                "type": "record",
                "name": "EpochStartInfo",
                "fields": [
                    {
                        "name": "TotalSupply",
                        "type": "string"
                    },
                    {
                        "name": "TotalToDistribute",
                        "type": "string"
                    },
                    {
                        "name": "TotalNewlyMinted",
                        "type": "string"
                    },
                    {
                        "name": "RewardsPerBlock",
                        "type": "string"
                    },
                    {
                        "name": "RewardsForProtocolSustainability",
                        "type": "string"
                    },
                    {
                        "name": "NodePrice",
                        "type": "string"
                    },
                    {
                        "name": "PrevEpochStartRound",
                        "type": "int"
                    },
                    {
                        "name": "PrevEpochStartHash",
                        "type": "string"
                    }
                ]
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _MiniBlock_schema, _MiniBlock_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "MiniBlock",
    "fields": [
        {
            "name": "Hash",
            "type": "string"
        },
        {
            "name": "SenderShardID",
            "type": "int"
        },
        {
            "name": "ReceiverShardID",
            "type": "int"
        },
        {
            "name": "SenderBlockHash",
            "type": "string"
        },
        {
            "name": "ReceiverBlockHash",
            "type": "string"
        },
        {
            "name": "Type",
            "type": "string"
        },
        {
            "name": "Timestamp",
            "type": "long"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _EpochStartInfo_schema, _EpochStartInfo_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "EpochStartInfo",
    "fields": [
        {
            "name": "TotalSupply",
            "type": "string"
        },
        {
            "name": "TotalToDistribute",
            "type": "string"
        },
        {
            "name": "TotalNewlyMinted",
            "type": "string"
        },
        {
            "name": "RewardsPerBlock",
            "type": "string"
        },
        {
            "name": "RewardsForProtocolSustainability",
            "type": "string"
        },
        {
            "name": "NodePrice",
            "type": "string"
        },
        {
            "name": "PrevEpochStartRound",
            "type": "int"
        },
        {
            "name": "PrevEpochStartHash",
            "type": "string"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Transaction_schema, _Transaction_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Transaction",
    "fields": [
        {
            "name": "Hash",
            "type": "string"
        },
        {
            "name": "MiniBlockHash",
            "type": "string"
        },
        {
            "name": "BlockHash",
            "type": "string"
        },
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "Round",
            "type": "long"
        },
        {
            "name": "Value",
            "type": "string"
        },
        {
            "name": "Receiver",
            "type": "string"
        },
        {
            "name": "Sender",
            "type": "string"
        },
        {
            "name": "ReceiverShard",
            "type": "int"
        },
        {
            "name": "SenderShard",
            "type": "int"
        },
        {
            "name": "GasPrice",
            "type": "long"
        },
        {
            "name": "GasLimit",
            "type": "long"
        },
        {
            "name": "Signature",
            "type": "string"
        },
        {
            "name": "Timestamp",
            "type": "long"
        },
        {
            "name": "SenderUserName",
            "type": "bytes"
        },
        {
            "name": "ReceiverUserName",
            "type": "bytes"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _SCResult_schema, _SCResult_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SCResult",
    "fields": [
        {
            "name": "Hash",
            "type": "string"
        },
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "GasLimit",
            "type": "long"
        },
        {
            "name": "GasPrice",
            "type": "long"
        },
        {
            "name": "Value",
            "type": "string"
        },
        {
            "name": "Sender",
            "type": "string"
        },
        {
            "name": "Receiver",
            "type": "string"
        },
        {
            "name": "RelayerAddr",
            "type": "string"
        },
        {
            "name": "RelayedValue",
            "type": "string"
        },
        {
            "name": "Code",
            "type": "string"
        },
        {
            "name": "Data",
            "type": "bytes"
        },
        {
            "name": "PrevTxHash",
            "type": "string"
        },
        {
            "name": "OriginalTxHash",
            "type": "string"
        },
        {
            "name": "CallType",
            "type": "string"
        },
        {
            "name": "CodeMetadata",
            "type": "bytes"
        },
        {
            "name": "ReturnMessage",
            "type": "string"
        },
        {
            "name": "Timestamp",
            "type": "long"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Receipt_schema, _Receipt_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Receipt",
    "fields": [
        {
            "name": "Hash",
            "type": "string"
        },
        {
            "name": "Value",
            "type": "string"
        },
        {
            "name": "Sender",
            "type": "string"
        },
        {
            "name": "Data",
            "type": "string"
        },
        {
            "name": "TxHash",
            "type": "string"
        },
        {
            "name": "Timestamp",
            "type": "long"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Log_schema, _Log_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Log",
    "fields": [
        {
            "name": "ID",
            "type": "string"
        },
        {
            "name": "Address",
            "type": "string"
        },
        {
            "name": "Events",
            "type": {
                "type": "array",
                "items": {
                    "type": "record",
                    "name": "Event",
                    "fields": [
                        {
                            "name": "Address",
                            "type": "string"
                        },
                        {
                            "name": "Identifier",
                            "type": "string"
                        },
                        {
                            "name": "Topics",
                            "type": {
                                "type": "array",
                                "items": "bytes"
                            }
                        },
                        {
                            "name": "Data",
                            "type": "bytes"
                        }
                    ]
                }
            }
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Event_schema, _Event_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Event",
    "fields": [
        {
            "name": "Address",
            "type": "string"
        },
        {
            "name": "Identifier",
            "type": "string"
        },
        {
            "name": "Topics",
            "type": {
                "type": "array",
                "items": "bytes"
            }
        },
        {
            "name": "Data",
            "type": "bytes"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _AccountBalanceUpdate_schema, _AccountBalanceUpdate_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "AccountBalanceUpdate",
    "fields": [
        {
            "name": "Address",
            "type": "string"
        },
        {
            "name": "Balance",
            "type": "string"
        },
        {
            "name": "Nonce",
            "type": "long"
        }
    ]
}`)
