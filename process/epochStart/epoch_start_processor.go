package epochStart

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-covalent-go/process/utility"
	"github.com/multiversx/mx-chain-covalent-go/schema"
	"github.com/multiversx/mx-chain-core-go/data/api"
)

type epochStartInfoProcessor struct {
}

// NewEpochStartInfoProcessor will create a new instance of epoch start info processor
func NewEpochStartInfoProcessor() *epochStartInfoProcessor {
	return &epochStartInfoProcessor{}
}

// ProcessEpochStartInfo converts receipts api epoch start info to a specific structure defined by avro schema
func (esi *epochStartInfoProcessor) ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schema.EpochStartInfo, error) {
	if apiEpochInfo == nil {
		return schema.NewEpochStartInfo(), nil
	}

	totalSupply, err := utility.GetBigIntBytesFromStr(apiEpochInfo.TotalSupply)
	if err != nil {
		return nil, err
	}
	totalToDistribute, err := utility.GetBigIntBytesFromStr(apiEpochInfo.TotalToDistribute)
	if err != nil {
		return nil, err
	}
	totalNewlyMinted, err := utility.GetBigIntBytesFromStr(apiEpochInfo.TotalNewlyMinted)
	if err != nil {
		return nil, err
	}
	rewardsPerBlock, err := utility.GetBigIntBytesFromStr(apiEpochInfo.RewardsPerBlock)
	if err != nil {
		return nil, err
	}
	rewardsForProtocolSustainability, err := utility.GetBigIntBytesFromStr(apiEpochInfo.RewardsForProtocolSustainability)
	if err != nil {
		return nil, err
	}
	nodePrice, err := utility.GetBigIntBytesFromStr(apiEpochInfo.NodePrice)
	if err != nil {
		return nil, err
	}
	prevEpochStartHash, err := hex.DecodeString(apiEpochInfo.PrevEpochStartHash)
	if err != nil {
		return nil, err
	}

	return &schema.EpochStartInfo{
		TotalSupply:                      totalSupply,
		TotalToDistribute:                totalToDistribute,
		TotalNewlyMinted:                 totalNewlyMinted,
		RewardsPerBlock:                  rewardsPerBlock,
		RewardsForProtocolSustainability: rewardsForProtocolSustainability,
		NodePrice:                        nodePrice,
		PrevEpochStartRound:              int64(apiEpochInfo.PrevEpochStartRound),
		PrevEpochStartHash:               prevEpochStartHash,
	}, nil
}
