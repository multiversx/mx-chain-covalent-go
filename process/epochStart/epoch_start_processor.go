package epochStart

import (
	"encoding/hex"

	"github.com/ElrondNetwork/covalent-indexer-go/process/utility"
	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
)

type epochStartInfoProcessor struct {
}

func NewEpochStartInfoProcessor() *epochStartInfoProcessor {
	return &epochStartInfoProcessor{}
}

func (esi *epochStartInfoProcessor) ProcessEpochStartInfo(apiEpochInfo *api.EpochStartInfo) (*schemaV2.EpochStartInfo, error) {
	if apiEpochInfo == nil {
		return schemaV2.NewEpochStartInfo(), nil
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

	return &schemaV2.EpochStartInfo{
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
