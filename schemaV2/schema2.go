package schemaV2

import "github.com/elodina/go-avro"

type HyperBlock struct {
	Hash                   []byte
	PrevBlockHash          []byte
	StateRootHash          []byte
	Nonce                  int64
	Round                  int64
	Epoch                  int32
	NumTxs                 int32
	AccumulatedFees        []byte
	DeveloperFees          []byte
	AccumulatedFeesInEpoch []byte
	DeveloperFeesInEpoch   []byte
	Timestamp              int64
	EpochStartInfo         *EpochStartInfo
	ShardBlocks            []*ShardBlocks
	Transactions           []*Transaction
	SCResults              []*SCResult
	StateChanges           []*AccountBalanceUpdate
	Status                 string
}

func NewHyperBlock() *HyperBlock {
	return &HyperBlock{
		Hash:                   make([]byte, 32),
		PrevBlockHash:          make([]byte, 32),
		StateRootHash:          make([]byte, 32),
		AccumulatedFees:        []byte{},
		DeveloperFees:          []byte{},
		AccumulatedFeesInEpoch: []byte{},
		DeveloperFeesInEpoch:   []byte{},
		StateChanges:           make([]*AccountBalanceUpdate, 0),
	}
}

func (o *HyperBlock) Schema() avro.Schema {
	if _HyperBlock_schema_err != nil {
		panic(_HyperBlock_schema_err)
	}
	return _HyperBlock_schema
}

type EpochStartInfo struct {
	TotalSupply                      []byte
	TotalToDistribute                []byte
	TotalNewlyMinted                 []byte
	RewardsPerBlock                  []byte
	RewardsForProtocolSustainability []byte
	NodePrice                        []byte
	PrevEpochStartRound              int64
	PrevEpochStartHash               []byte
}

func NewEpochStartInfo() *EpochStartInfo {
	return &EpochStartInfo{
		TotalSupply:                      []byte{},
		TotalToDistribute:                []byte{},
		TotalNewlyMinted:                 []byte{},
		RewardsPerBlock:                  []byte{},
		RewardsForProtocolSustainability: []byte{},
		NodePrice:                        []byte{},
	}
}

func (o *EpochStartInfo) Schema() avro.Schema {
	if _EpochStartInfo_schema_err != nil {
		panic(_EpochStartInfo_schema_err)
	}
	return _EpochStartInfo_schema
}

type ShardBlocks struct {
	Hash  []byte
	Nonce int64
	Round int64
	Shard int32
}

func NewShardBlocks() *ShardBlocks {
	return &ShardBlocks{
		Hash: make([]byte, 32),
	}
}

func (o *ShardBlocks) Schema() avro.Schema {
	if _ShardBlocks_schema_err != nil {
		panic(_ShardBlocks_schema_err)
	}
	return _ShardBlocks_schema
}

type Transaction struct {
	Type                              string
	ProcessingTypeOnSource            string
	ProcessingTypeOnDestination       string
	Hash                              []byte
	Nonce                             int64
	Round                             int64
	Epoch                             int32
	Value                             []byte
	Receiver                          []byte
	Sender                            []byte
	SenderUserName                    []byte
	ReceiverUserName                  []byte
	GasPrice                          int64
	GasLimit                          int64
	Data                              []byte
	CodeMetadata                      []byte
	Code                              []byte
	PreviousTransactionHash           []byte
	OriginalTransactionHash           []byte
	ReturnMessage                     string
	OriginalSender                    []byte
	Signature                         []byte
	SourceShard                       int32
	DestinationShard                  int32
	BlockNonce                        int64
	BlockHash                         []byte
	NotarizedAtSourceInMetaNonce      int64
	NotarizedAtSourceInMetaHash       []byte
	NotarizedAtDestinationInMetaNonce int64
	NotarizedAtDestinationInMetaHash  []byte
	MiniBlockType                     string
	MiniBlockHash                     []byte
	HyperBlockNonce                   int64
	HyperBlockHash                    []byte
	Timestamp                         int64
	Receipt                           *Receipt
	Log                               *Log
	Status                            string
	Tokens                            []string
	ESDTValues                        [][]byte
	Receivers                         [][]byte
	ReceiversShardIDs                 []int32
	Operation                         string
	Function                          string
	InitiallyPaidFee                  []byte
	IsRelayed                         bool
	IsRefund                          bool
}

func NewTransaction() *Transaction {
	return &Transaction{
		Hash:                             make([]byte, 32),
		Value:                            []byte{},
		Receiver:                         make([]byte, 62),
		Sender:                           make([]byte, 62),
		SenderUserName:                   []byte{},
		ReceiverUserName:                 []byte{},
		Data:                             []byte{},
		CodeMetadata:                     []byte{},
		Code:                             []byte{},
		PreviousTransactionHash:          make([]byte, 32),
		OriginalTransactionHash:          make([]byte, 32),
		OriginalSender:                   make([]byte, 62),
		BlockHash:                        make([]byte, 32),
		NotarizedAtSourceInMetaHash:      make([]byte, 32),
		NotarizedAtDestinationInMetaHash: make([]byte, 32),
		MiniBlockHash:                    make([]byte, 32),
		HyperBlockHash:                   make([]byte, 32),
		Receipt:                          NewReceipt(),
		Log:                              NewLog(),
		Tokens:                           make([]string, 0),
		ESDTValues:                       make([][]byte, 0),
		Receivers:                        make([][]byte, 0),
		ReceiversShardIDs:                make([]int32, 0),
		InitiallyPaidFee:                 []byte{},
	}
}

func (o *Transaction) Schema() avro.Schema {
	if _Transaction_schema_err != nil {
		panic(_Transaction_schema_err)
	}
	return _Transaction_schema
}

type Receipt struct {
	TxHash []byte
	Value  []byte
	Sender []byte
	Data   []byte
}

func NewReceipt() *Receipt {
	return &Receipt{
		TxHash: make([]byte, 32),
		Value:  []byte{},
		Sender: make([]byte, 62),
		Data:   []byte{},
	}
}

func (o *Receipt) Schema() avro.Schema {
	if _Receipt_schema_err != nil {
		panic(_Receipt_schema_err)
	}
	return _Receipt_schema
}

type Log struct {
	Address []byte
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
	Address    []byte
	Identifier []byte
	Topics     [][]byte
	Data       []byte
}

func NewEvent() *Event {
	return &Event{
		Identifier: []byte{},
		Topics:     make([][]byte, 0),
		Data:       []byte{},
	}
}

func (o *Event) Schema() avro.Schema {
	if _Event_schema_err != nil {
		panic(_Event_schema_err)
	}
	return _Event_schema
}

type SCResult struct {
	Hash              []byte
	Nonce             int64
	Value             []byte
	Receiver          []byte
	Sender            []byte
	Relayer           []byte
	RelayedValue      []byte
	Code              []byte
	Data              []byte
	PrevTxHash        []byte
	OriginalTxHash    []byte
	GasLimit          int64
	GasPrice          int64
	CallType          int32
	CodeMetadata      []byte
	ReturnMessage     []byte
	OriginalSender    []byte
	Log               *Log
	Tokens            []string
	ESDTValues        [][]byte
	Receivers         [][]byte
	ReceiversShardIDs []int32
	Operation         string
	Function          string
	IsRelayed         bool
	IsRefund          bool
}

func NewSCResult() *SCResult {
	return &SCResult{
		Hash:              make([]byte, 32),
		Value:             []byte{},
		Receiver:          make([]byte, 62),
		Sender:            make([]byte, 62),
		RelayedValue:      []byte{},
		Code:              []byte{},
		Data:              []byte{},
		PrevTxHash:        make([]byte, 32),
		OriginalTxHash:    make([]byte, 32),
		CodeMetadata:      []byte{},
		ReturnMessage:     []byte{},
		OriginalSender:    make([]byte, 62),
		Tokens:            make([]string, 0),
		ESDTValues:        make([][]byte, 0),
		Receivers:         make([][]byte, 0),
		ReceiversShardIDs: make([]int32, 0),
	}
}

func (o *SCResult) Schema() avro.Schema {
	if _SCResult_schema_err != nil {
		panic(_SCResult_schema_err)
	}
	return _SCResult_schema
}

type AccountBalanceUpdate struct {
	Address []byte
	Balance []byte
	Nonce   int64
}

func NewAccountBalanceUpdate() *AccountBalanceUpdate {
	return &AccountBalanceUpdate{
		Address: make([]byte, 62),
		Balance: []byte{},
	}
}

func (o *AccountBalanceUpdate) Schema() avro.Schema {
	if _AccountBalanceUpdate_schema_err != nil {
		panic(_AccountBalanceUpdate_schema_err)
	}
	return _AccountBalanceUpdate_schema
}

// Generated by codegen. Please do not modify.
var _HyperBlock_schema, _HyperBlock_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "com.covalenthq.block.schema",
    "name": "HyperBlock",
    "fields": [
        {
            "name": "Hash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "PrevBlockHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "StateRootHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
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
            "name": "Epoch",
            "type": "int"
        },
        {
            "name": "NumTxs",
            "type": "int"
        },
        {
            "name": "AccumulatedFees",
            "type": "bytes"
        },
        {
            "name": "DeveloperFees",
            "type": "bytes"
        },
        {
            "name": "AccumulatedFeesInEpoch",
            "type": "bytes"
        },
        {
            "name": "DeveloperFeesInEpoch",
            "type": "bytes"
        },
        {
            "name": "Timestamp",
            "type": "long"
        },
        {
            "name": "EpochStartInfo",
            "default": null,
            "type": [
                "null",
                {
                    "type": "record",
                    "name": "EpochStartInfo",
                    "fields": [
                        {
                            "name": "TotalSupply",
                            "type": "bytes"
                        },
                        {
                            "name": "TotalToDistribute",
                            "type": "bytes"
                        },
                        {
                            "name": "TotalNewlyMinted",
                            "type": "bytes"
                        },
                        {
                            "name": "RewardsPerBlock",
                            "type": "bytes"
                        },
                        {
                            "name": "RewardsForProtocolSustainability",
                            "type": "bytes"
                        },
                        {
                            "name": "NodePrice",
                            "type": "bytes"
                        },
                        {
                            "name": "PrevEpochStartRound",
                            "type": "long"
                        },
                        {
                            "name": "PrevEpochStartHash",
                            "default": null,
                            "type": [
                                "null",
                                {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        {
            "name": "ShardBlocks",
            "default": null,
            "type": [
                "null",
                {
                    "type": "array",
                    "items": {
                        "type": "record",
                        "name": "ShardBlocks",
                        "fields": [
                            {
                                "name": "Hash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
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
                                "name": "Shard",
                                "type": "int"
                            }
                        ]
                    }
                }
            ]
        },
        {
            "name": "Transactions",
            "default": null,
            "type": [
                "null",
                {
                    "type": "array",
                    "items": {
                        "type": "record",
                        "name": "Transaction",
                        "fields": [
                            {
                                "name": "Type",
                                "type": "string"
                            },
                            {
                                "name": "ProcessingTypeOnSource",
                                "type": "string"
                            },
                            {
                                "name": "ProcessingTypeOnDestination",
                                "type": "string"
                            },
                            {
                                "name": "Hash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
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
                                "name": "Epoch",
                                "type": "int"
                            },
                            {
                                "name": "Value",
                                "type": "bytes"
                            },
                            {
                                "name": "Receiver",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "Sender",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "SenderUserName",
                                "type": "bytes"
                            },
                            {
                                "name": "ReceiverUserName",
                                "type": "bytes"
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
                                "name": "Data",
                                "type": "bytes"
                            },
                            {
                                "name": "CodeMetadata",
                                "type": "bytes"
                            },
                            {
                                "name": "Code",
                                "type": "bytes"
                            },
                            {
                                "name": "PreviousTransactionHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "OriginalTransactionHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "ReturnMessage",
                                "type": "string"
                            },
                            {
                                "name": "OriginalSender",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "Signature",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 64,
                                        "name": "signature"
                                    }
                                ]
                            },
                            {
                                "name": "SourceShard",
                                "type": "int"
                            },
                            {
                                "name": "DestinationShard",
                                "type": "int"
                            },
                            {
                                "name": "BlockNonce",
                                "type": "long"
                            },
                            {
                                "name": "BlockHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "NotarizedAtSourceInMetaNonce",
                                "type": "long"
                            },
                            {
                                "name": "NotarizedAtSourceInMetaHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "NotarizedAtDestinationInMetaNonce",
                                "type": "long"
                            },
                            {
                                "name": "NotarizedAtDestinationInMetaHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "MiniBlockType",
                                "type": "string"
                            },
                            {
                                "name": "MiniBlockHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "HyperBlockNonce",
                                "type": "long"
                            },
                            {
                                "name": "HyperBlockHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "Timestamp",
                                "type": "long"
                            },
                            {
                                "name": "Receipt",
                                "type": {
                                    "type": "record",
                                    "name": "Receipt",
                                    "fields": [
                                        {
                                            "name": "TxHash",
                                            "type": {
                                                "type": "fixed",
                                                "size": 32,
                                                "name": "hash"
                                            }
                                        },
                                        {
                                            "name": "Value",
                                            "type": "bytes"
                                        },
                                        {
                                            "name": "Sender",
                                            "type": {
                                                "type": "fixed",
                                                "size": 62,
                                                "name": "address"
                                            }
                                        },
                                        {
                                            "name": "Data",
                                            "type": "bytes"
                                        }
                                    ]
                                }
                            },
                            {
                                "name": "Log",
                                "type": {
                                    "type": "record",
                                    "name": "Log",
                                    "fields": [
                                        {
                                            "name": "Address",
                                            "default": null,
                                            "type": [
                                                "null",
                                                {
                                                    "type": "fixed",
                                                    "size": 62,
                                                    "name": "address"
                                                }
                                            ]
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
                                                            "default": null,
                                                            "type": [
                                                                "null",
                                                                {
                                                                    "type": "fixed",
                                                                    "size": 62,
                                                                    "name": "address"
                                                                }
                                                            ]
                                                        },
                                                        {
                                                            "name": "Identifier",
                                                            "type": "bytes"
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
                            },
                            {
                                "name": "Status",
                                "type": "string"
                            },
                            {
                                "name": "Tokens",
                                "type": {
                                    "type": "array",
                                    "items": "string"
                                }
                            },
                            {
                                "name": "ESDTValues",
                                "type": {
                                    "type": "array",
                                    "items": "bytes"
                                }
                            },
                            {
                                "name": "Receivers",
                                "type": {
                                    "type": "array",
                                    "items": {
                                        "type": "fixed",
                                        "size": 62,
                                        "name": "address"
                                    }
                                }
                            },
                            {
                                "name": "ReceiversShardIDs",
                                "type": {
                                    "type": "array",
                                    "items": "int"
                                }
                            },
                            {
                                "name": "Operation",
                                "type": "string"
                            },
                            {
                                "name": "Function",
                                "type": "string"
                            },
                            {
                                "name": "InitiallyPaidFee",
                                "type": "bytes"
                            },
                            {
                                "name": "IsRelayed",
                                "type": "boolean"
                            },
                            {
                                "name": "IsRefund",
                                "type": "boolean"
                            }
                        ]
                    }
                }
            ]
        },
        {
            "name": "SCResults",
            "default": null,
            "type": [
                "null",
                {
                    "type": "array",
                    "items": {
                        "type": "record",
                        "name": "SCResult",
                        "fields": [
                            {
                                "name": "Hash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "Nonce",
                                "type": "long"
                            },
                            {
                                "name": "Value",
                                "type": "bytes"
                            },
                            {
                                "name": "Receiver",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "Sender",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "Relayer",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 62,
                                        "name": "address"
                                    }
                                ]
                            },
                            {
                                "name": "RelayedValue",
                                "type": "bytes"
                            },
                            {
                                "name": "Code",
                                "type": "bytes"
                            },
                            {
                                "name": "Data",
                                "type": "bytes"
                            },
                            {
                                "name": "PrevTxHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
                            },
                            {
                                "name": "OriginalTxHash",
                                "type": {
                                    "type": "fixed",
                                    "size": 32,
                                    "name": "hash"
                                }
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
                                "name": "CallType",
                                "type": "int"
                            },
                            {
                                "name": "CodeMetadata",
                                "type": "bytes"
                            },
                            {
                                "name": "ReturnMessage",
                                "type": "bytes"
                            },
                            {
                                "name": "OriginalSender",
                                "type": {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            },
                            {
                                "name": "Log",
                                "type": "Log"
                            },
                            {
                                "name": "Tokens",
                                "type": {
                                    "type": "array",
                                    "items": "string"
                                }
                            },
                            {
                                "name": "ESDTValues",
                                "type": {
                                    "type": "array",
                                    "items": "bytes"
                                }
                            },
                            {
                                "name": "Receivers",
                                "type": {
                                    "type": "array",
                                    "items": {
                                        "type": "fixed",
                                        "size": 62,
                                        "name": "address"
                                    }
                                }
                            },
                            {
                                "name": "ReceiversShardIDs",
                                "type": {
                                    "type": "array",
                                    "items": "int"
                                }
                            },
                            {
                                "name": "Operation",
                                "type": "string"
                            },
                            {
                                "name": "Function",
                                "type": "string"
                            },
                            {
                                "name": "IsRelayed",
                                "type": "boolean"
                            },
                            {
                                "name": "IsRefund",
                                "type": "boolean"
                            }
                        ]
                    }
                }
            ]
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
                            "type": {
                                "type": "fixed",
                                "size": 62,
                                "name": "address"
                            }
                        },
                        {
                            "name": "Balance",
                            "type": "bytes"
                        },
                        {
                            "name": "Nonce",
                            "type": "long"
                        }
                    ]
                }
            }
        },
        {
            "name": "Status",
            "type": "string"
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
            "type": "bytes"
        },
        {
            "name": "TotalToDistribute",
            "type": "bytes"
        },
        {
            "name": "TotalNewlyMinted",
            "type": "bytes"
        },
        {
            "name": "RewardsPerBlock",
            "type": "bytes"
        },
        {
            "name": "RewardsForProtocolSustainability",
            "type": "bytes"
        },
        {
            "name": "NodePrice",
            "type": "bytes"
        },
        {
            "name": "PrevEpochStartRound",
            "type": "long"
        },
        {
            "name": "PrevEpochStartHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _ShardBlocks_schema, _ShardBlocks_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "ShardBlocks",
    "fields": [
        {
            "name": "Hash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
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
            "name": "Shard",
            "type": "int"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Transaction_schema, _Transaction_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Transaction",
    "fields": [
        {
            "name": "Type",
            "type": "string"
        },
        {
            "name": "ProcessingTypeOnSource",
            "type": "string"
        },
        {
            "name": "ProcessingTypeOnDestination",
            "type": "string"
        },
        {
            "name": "Hash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
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
            "name": "Epoch",
            "type": "int"
        },
        {
            "name": "Value",
            "type": "bytes"
        },
        {
            "name": "Receiver",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Sender",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "SenderUserName",
            "type": "bytes"
        },
        {
            "name": "ReceiverUserName",
            "type": "bytes"
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
            "name": "Data",
            "type": "bytes"
        },
        {
            "name": "CodeMetadata",
            "type": "bytes"
        },
        {
            "name": "Code",
            "type": "bytes"
        },
        {
            "name": "PreviousTransactionHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "OriginalTransactionHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "ReturnMessage",
            "type": "string"
        },
        {
            "name": "OriginalSender",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Signature",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 64,
                    "name": "signature"
                }
            ]
        },
        {
            "name": "SourceShard",
            "type": "int"
        },
        {
            "name": "DestinationShard",
            "type": "int"
        },
        {
            "name": "BlockNonce",
            "type": "long"
        },
        {
            "name": "BlockHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "NotarizedAtSourceInMetaNonce",
            "type": "long"
        },
        {
            "name": "NotarizedAtSourceInMetaHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "NotarizedAtDestinationInMetaNonce",
            "type": "long"
        },
        {
            "name": "NotarizedAtDestinationInMetaHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "MiniBlockType",
            "type": "string"
        },
        {
            "name": "MiniBlockHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "HyperBlockNonce",
            "type": "long"
        },
        {
            "name": "HyperBlockHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "Timestamp",
            "type": "long"
        },
        {
            "name": "Receipt",
            "type": {
                "type": "record",
                "name": "Receipt",
                "fields": [
                    {
                        "name": "TxHash",
                        "type": {
                            "type": "fixed",
                            "size": 32,
                            "name": "hash"
                        }
                    },
                    {
                        "name": "Value",
                        "type": "bytes"
                    },
                    {
                        "name": "Sender",
                        "type": {
                            "type": "fixed",
                            "size": 62,
                            "name": "address"
                        }
                    },
                    {
                        "name": "Data",
                        "type": "bytes"
                    }
                ]
            }
        },
        {
            "name": "Log",
            "type": {
                "type": "record",
                "name": "Log",
                "fields": [
                    {
                        "name": "Address",
                        "default": null,
                        "type": [
                            "null",
                            {
                                "type": "fixed",
                                "size": 62,
                                "name": "address"
                            }
                        ]
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
                                        "default": null,
                                        "type": [
                                            "null",
                                            {
                                                "type": "fixed",
                                                "size": 62,
                                                "name": "address"
                                            }
                                        ]
                                    },
                                    {
                                        "name": "Identifier",
                                        "type": "bytes"
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
        },
        {
            "name": "Status",
            "type": "string"
        },
        {
            "name": "Tokens",
            "type": {
                "type": "array",
                "items": "string"
            }
        },
        {
            "name": "ESDTValues",
            "type": {
                "type": "array",
                "items": "bytes"
            }
        },
        {
            "name": "Receivers",
            "type": {
                "type": "array",
                "items": {
                    "type": "fixed",
                    "size": 62,
                    "name": "address"
                }
            }
        },
        {
            "name": "ReceiversShardIDs",
            "type": {
                "type": "array",
                "items": "int"
            }
        },
        {
            "name": "Operation",
            "type": "string"
        },
        {
            "name": "Function",
            "type": "string"
        },
        {
            "name": "InitiallyPaidFee",
            "type": "bytes"
        },
        {
            "name": "IsRelayed",
            "type": "boolean"
        },
        {
            "name": "IsRefund",
            "type": "boolean"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Receipt_schema, _Receipt_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Receipt",
    "fields": [
        {
            "name": "TxHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "Value",
            "type": "bytes"
        },
        {
            "name": "Sender",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Data",
            "type": "bytes"
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _Log_schema, _Log_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "Log",
    "fields": [
        {
            "name": "Address",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 62,
                    "name": "address"
                }
            ]
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
                            "default": null,
                            "type": [
                                "null",
                                {
                                    "type": "fixed",
                                    "size": 62,
                                    "name": "address"
                                }
                            ]
                        },
                        {
                            "name": "Identifier",
                            "type": "bytes"
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
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 62,
                    "name": "address"
                }
            ]
        },
        {
            "name": "Identifier",
            "type": "bytes"
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
var _SCResult_schema, _SCResult_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "SCResult",
    "fields": [
        {
            "name": "Hash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "Value",
            "type": "bytes"
        },
        {
            "name": "Receiver",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Sender",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Relayer",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 62,
                    "name": "address"
                }
            ]
        },
        {
            "name": "RelayedValue",
            "type": "bytes"
        },
        {
            "name": "Code",
            "type": "bytes"
        },
        {
            "name": "Data",
            "type": "bytes"
        },
        {
            "name": "PrevTxHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
        },
        {
            "name": "OriginalTxHash",
            "type": {
                "type": "fixed",
                "size": 32,
                "name": "hash"
            }
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
            "name": "CallType",
            "type": "int"
        },
        {
            "name": "CodeMetadata",
            "type": "bytes"
        },
        {
            "name": "ReturnMessage",
            "type": "bytes"
        },
        {
            "name": "OriginalSender",
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Log",
            "type": "Log"
        },
        {
            "name": "Tokens",
            "type": {
                "type": "array",
                "items": "string"
            }
        },
        {
            "name": "ESDTValues",
            "type": {
                "type": "array",
                "items": "bytes"
            }
        },
        {
            "name": "Receivers",
            "type": {
                "type": "array",
                "items": {
                    "type": "fixed",
                    "size": 62,
                    "name": "address"
                }
            }
        },
        {
            "name": "ReceiversShardIDs",
            "type": {
                "type": "array",
                "items": "int"
            }
        },
        {
            "name": "Operation",
            "type": "string"
        },
        {
            "name": "Function",
            "type": "string"
        },
        {
            "name": "IsRelayed",
            "type": "boolean"
        },
        {
            "name": "IsRefund",
            "type": "boolean"
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
            "type": {
                "type": "fixed",
                "size": 62,
                "name": "address"
            }
        },
        {
            "name": "Balance",
            "type": "bytes"
        },
        {
            "name": "Nonce",
            "type": "long"
        }
    ]
}`)
