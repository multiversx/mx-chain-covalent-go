package schema

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
	Status                 string
}

func NewHyperBlock() *HyperBlock {
	return &HyperBlock{
		Hash:                   make([]byte, 32),
		AccumulatedFees:        []byte{},
		DeveloperFees:          []byte{},
		AccumulatedFeesInEpoch: []byte{},
		DeveloperFeesInEpoch:   []byte{},
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
	Hash            []byte
	Nonce           int64
	Round           int64
	Shard           int32
	RootHash        []byte
	MiniBlockHashes [][]byte
	StateChanges    []*AccountBalanceUpdate
}

func NewShardBlocks() *ShardBlocks {
	return &ShardBlocks{
		Hash:         make([]byte, 32),
		StateChanges: make([]*AccountBalanceUpdate, 0),
	}
}

func (o *ShardBlocks) Schema() avro.Schema {
	if _ShardBlocks_schema_err != nil {
		panic(_ShardBlocks_schema_err)
	}
	return _ShardBlocks_schema
}

type AccountBalanceUpdate struct {
	Address []byte
	Balance []byte
	Nonce   int64
	Tokens  []*AccountTokenData
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

type AccountTokenData struct {
	Nonce      int64
	Identifier string
	Balance    []byte
	Properties string
	MetaData   *MetaData
}

func NewAccountTokenData() *AccountTokenData {
	return &AccountTokenData{
		Balance: []byte{},
	}
}

func (o *AccountTokenData) Schema() avro.Schema {
	if _AccountTokenData_schema_err != nil {
		panic(_AccountTokenData_schema_err)
	}
	return _AccountTokenData_schema
}

type MetaData struct {
	Nonce      int64
	Name       []byte
	Creator    []byte
	Royalties  int32
	Hash       []byte
	URIs       [][]byte
	Attributes []byte
}

func NewMetaData() *MetaData {
	return &MetaData{
		Name:       []byte{},
		Creator:    []byte{},
		Hash:       []byte{},
		URIs:       make([][]byte, 0),
		Attributes: []byte{},
	}
}

func (o *MetaData) Schema() avro.Schema {
	if _MetaData_schema_err != nil {
		panic(_MetaData_schema_err)
	}
	return _MetaData_schema
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
	CallType                          string
	RelayerAddress                    []byte
	RelayedValue                      []byte
	ChainID                           string
	Version                           int32
	Options                           int32
}

func NewTransaction() *Transaction {
	return &Transaction{
		Hash:              make([]byte, 32),
		Value:             []byte{},
		Receiver:          make([]byte, 62),
		Sender:            make([]byte, 62),
		SenderUserName:    []byte{},
		ReceiverUserName:  []byte{},
		Data:              []byte{},
		CodeMetadata:      []byte{},
		Code:              []byte{},
		MiniBlockHash:     make([]byte, 32),
		Tokens:            make([]string, 0),
		ESDTValues:        make([][]byte, 0),
		Receivers:         make([][]byte, 0),
		ReceiversShardIDs: make([]int32, 0),
		InitiallyPaidFee:  []byte{},
		RelayedValue:      []byte{},
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
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "StateRootHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
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
                            },
                            {
                                "name": "RootHash",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "MiniBlockHashes",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "array",
                                        "items": {
                                            "type": "fixed",
                                            "size": 32,
                                            "name": "hash"
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
                                            },
                                            {
                                                "name": "Tokens",
                                                "default": null,
                                                "type": [
                                                    "null",
                                                    {
                                                        "type": "array",
                                                        "items": {
                                                            "type": "record",
                                                            "name": "AccountTokenData",
                                                            "fields": [
                                                                {
                                                                    "name": "Nonce",
                                                                    "type": "long"
                                                                },
                                                                {
                                                                    "name": "Identifier",
                                                                    "type": "string"
                                                                },
                                                                {
                                                                    "name": "Balance",
                                                                    "type": "bytes"
                                                                },
                                                                {
                                                                    "name": "Properties",
                                                                    "type": "string"
                                                                },
                                                                {
                                                                    "name": "MetaData",
                                                                    "default": null,
                                                                    "type": [
                                                                        "null",
                                                                        {
                                                                            "type": "record",
                                                                            "name": "MetaData",
                                                                            "fields": [
                                                                                {
                                                                                    "name": "Nonce",
                                                                                    "type": "long"
                                                                                },
                                                                                {
                                                                                    "name": "Name",
                                                                                    "type": "bytes"
                                                                                },
                                                                                {
                                                                                    "name": "Creator",
                                                                                    "type": "bytes"
                                                                                },
                                                                                {
                                                                                    "name": "Royalties",
                                                                                    "type": "int"
                                                                                },
                                                                                {
                                                                                    "name": "Hash",
                                                                                    "type": "bytes"
                                                                                },
                                                                                {
                                                                                    "name": "URIs",
                                                                                    "type": {
                                                                                        "type": "array",
                                                                                        "items": "bytes"
                                                                                    }
                                                                                },
                                                                                {
                                                                                    "name": "Attributes",
                                                                                    "type": "bytes"
                                                                                }
                                                                            ]
                                                                        }
                                                                    ]
                                                                }
                                                            ]
                                                        }
                                                    }
                                                ]
                                            }
                                        ]
                                    }
                                }
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
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "OriginalTransactionHash",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "ReturnMessage",
                                "type": "string"
                            },
                            {
                                "name": "OriginalSender",
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
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "NotarizedAtSourceInMetaNonce",
                                "type": "long"
                            },
                            {
                                "name": "NotarizedAtSourceInMetaHash",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "NotarizedAtDestinationInMetaNonce",
                                "type": "long"
                            },
                            {
                                "name": "NotarizedAtDestinationInMetaHash",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
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
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "fixed",
                                        "size": 32,
                                        "name": "hash"
                                    }
                                ]
                            },
                            {
                                "name": "Timestamp",
                                "type": "long"
                            },
                            {
                                "name": "Receipt",
                                "default": null,
                                "type": [
                                    "null",
                                    {
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
                                ]
                            },
                            {
                                "name": "Log",
                                "default": null,
                                "type": [
                                    "null",
                                    {
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
                                ]
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
                            },
                            {
                                "name": "CallType",
                                "type": "string"
                            },
                            {
                                "name": "RelayerAddress",
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
                                "name": "ChainID",
                                "type": "string"
                            },
                            {
                                "name": "Version",
                                "type": "int"
                            },
                            {
                                "name": "Options",
                                "type": "int"
                            }
                        ]
                    }
                }
            ]
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
        },
        {
            "name": "RootHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "MiniBlockHashes",
            "default": null,
            "type": [
                "null",
                {
                    "type": "array",
                    "items": {
                        "type": "fixed",
                        "size": 32,
                        "name": "hash"
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
                        },
                        {
                            "name": "Tokens",
                            "default": null,
                            "type": [
                                "null",
                                {
                                    "type": "array",
                                    "items": {
                                        "type": "record",
                                        "name": "AccountTokenData",
                                        "fields": [
                                            {
                                                "name": "Nonce",
                                                "type": "long"
                                            },
                                            {
                                                "name": "Identifier",
                                                "type": "string"
                                            },
                                            {
                                                "name": "Balance",
                                                "type": "bytes"
                                            },
                                            {
                                                "name": "Properties",
                                                "type": "string"
                                            },
                                            {
                                                "name": "MetaData",
                                                "default": null,
                                                "type": [
                                                    "null",
                                                    {
                                                        "type": "record",
                                                        "name": "MetaData",
                                                        "fields": [
                                                            {
                                                                "name": "Nonce",
                                                                "type": "long"
                                                            },
                                                            {
                                                                "name": "Name",
                                                                "type": "bytes"
                                                            },
                                                            {
                                                                "name": "Creator",
                                                                "type": "bytes"
                                                            },
                                                            {
                                                                "name": "Royalties",
                                                                "type": "int"
                                                            },
                                                            {
                                                                "name": "Hash",
                                                                "type": "bytes"
                                                            },
                                                            {
                                                                "name": "URIs",
                                                                "type": {
                                                                    "type": "array",
                                                                    "items": "bytes"
                                                                }
                                                            },
                                                            {
                                                                "name": "Attributes",
                                                                "type": "bytes"
                                                            }
                                                        ]
                                                    }
                                                ]
                                            }
                                        ]
                                    }
                                }
                            ]
                        }
                    ]
                }
            }
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
        },
        {
            "name": "Tokens",
            "default": null,
            "type": [
                "null",
                {
                    "type": "array",
                    "items": {
                        "type": "record",
                        "name": "AccountTokenData",
                        "fields": [
                            {
                                "name": "Nonce",
                                "type": "long"
                            },
                            {
                                "name": "Identifier",
                                "type": "string"
                            },
                            {
                                "name": "Balance",
                                "type": "bytes"
                            },
                            {
                                "name": "Properties",
                                "type": "string"
                            },
                            {
                                "name": "MetaData",
                                "default": null,
                                "type": [
                                    "null",
                                    {
                                        "type": "record",
                                        "name": "MetaData",
                                        "fields": [
                                            {
                                                "name": "Nonce",
                                                "type": "long"
                                            },
                                            {
                                                "name": "Name",
                                                "type": "bytes"
                                            },
                                            {
                                                "name": "Creator",
                                                "type": "bytes"
                                            },
                                            {
                                                "name": "Royalties",
                                                "type": "int"
                                            },
                                            {
                                                "name": "Hash",
                                                "type": "bytes"
                                            },
                                            {
                                                "name": "URIs",
                                                "type": {
                                                    "type": "array",
                                                    "items": "bytes"
                                                }
                                            },
                                            {
                                                "name": "Attributes",
                                                "type": "bytes"
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                }
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _AccountTokenData_schema, _AccountTokenData_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "AccountTokenData",
    "fields": [
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "Identifier",
            "type": "string"
        },
        {
            "name": "Balance",
            "type": "bytes"
        },
        {
            "name": "Properties",
            "type": "string"
        },
        {
            "name": "MetaData",
            "default": null,
            "type": [
                "null",
                {
                    "type": "record",
                    "name": "MetaData",
                    "fields": [
                        {
                            "name": "Nonce",
                            "type": "long"
                        },
                        {
                            "name": "Name",
                            "type": "bytes"
                        },
                        {
                            "name": "Creator",
                            "type": "bytes"
                        },
                        {
                            "name": "Royalties",
                            "type": "int"
                        },
                        {
                            "name": "Hash",
                            "type": "bytes"
                        },
                        {
                            "name": "URIs",
                            "type": {
                                "type": "array",
                                "items": "bytes"
                            }
                        },
                        {
                            "name": "Attributes",
                            "type": "bytes"
                        }
                    ]
                }
            ]
        }
    ]
}`)

// Generated by codegen. Please do not modify.
var _MetaData_schema, _MetaData_schema_err = avro.ParseSchema(`{
    "type": "record",
    "name": "MetaData",
    "fields": [
        {
            "name": "Nonce",
            "type": "long"
        },
        {
            "name": "Name",
            "type": "bytes"
        },
        {
            "name": "Creator",
            "type": "bytes"
        },
        {
            "name": "Royalties",
            "type": "int"
        },
        {
            "name": "Hash",
            "type": "bytes"
        },
        {
            "name": "URIs",
            "type": {
                "type": "array",
                "items": "bytes"
            }
        },
        {
            "name": "Attributes",
            "type": "bytes"
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
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "OriginalTransactionHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "ReturnMessage",
            "type": "string"
        },
        {
            "name": "OriginalSender",
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
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "NotarizedAtSourceInMetaNonce",
            "type": "long"
        },
        {
            "name": "NotarizedAtSourceInMetaHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "NotarizedAtDestinationInMetaNonce",
            "type": "long"
        },
        {
            "name": "NotarizedAtDestinationInMetaHash",
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
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
            "default": null,
            "type": [
                "null",
                {
                    "type": "fixed",
                    "size": 32,
                    "name": "hash"
                }
            ]
        },
        {
            "name": "Timestamp",
            "type": "long"
        },
        {
            "name": "Receipt",
            "default": null,
            "type": [
                "null",
                {
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
            ]
        },
        {
            "name": "Log",
            "default": null,
            "type": [
                "null",
                {
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
            ]
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
        },
        {
            "name": "CallType",
            "type": "string"
        },
        {
            "name": "RelayerAddress",
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
            "name": "ChainID",
            "type": "string"
        },
        {
            "name": "Version",
            "type": "int"
        },
        {
            "name": "Options",
            "type": "int"
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
