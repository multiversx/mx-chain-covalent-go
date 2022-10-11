package epochStart

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ElrondNetwork/covalent-indexer-go/schemaV2"
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/stretchr/testify/require"
)

func TestEpochStartInfoProcessor_ProcessEpochStartInfo(t *testing.T) {
	t.Parallel()

	esi := NewEpochStartInfoProcessor()

	apiEpochStartInfo := &api.EpochStartInfo{
		TotalSupply:                      "4444",
		TotalToDistribute:                "3333",
		TotalNewlyMinted:                 "2222",
		RewardsPerBlock:                  "1111",
		RewardsForProtocolSustainability: "4321",
		NodePrice:                        "1234",
		PrevEpochStartRound:              4,
		PrevEpochStartHash:               "0f",
	}

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		epochStartInfo, err := esi.ProcessEpochStartInfo(apiEpochStartInfo)
		require.Nil(t, err)
		require.Equal(t, &schemaV2.EpochStartInfo{
			TotalSupply:                      big.NewInt(4444).Bytes(),
			TotalToDistribute:                big.NewInt(3333).Bytes(),
			TotalNewlyMinted:                 big.NewInt(2222).Bytes(),
			RewardsPerBlock:                  big.NewInt(1111).Bytes(),
			RewardsForProtocolSustainability: big.NewInt(4321).Bytes(),
			NodePrice:                        big.NewInt(1234).Bytes(),
			PrevEpochStartRound:              4,
			PrevEpochStartHash:               []byte{15},
		}, epochStartInfo)
	})

	t.Run("nil api epoch info, should return nil info", func(t *testing.T) {
		t.Parallel()

		epochStartInfo, err := esi.ProcessEpochStartInfo(nil)
		require.Nil(t, err)
		require.Nil(t, epochStartInfo)
	})

	t.Run("invalid total supply, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.TotalSupply = "total supply"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "total supply"))
	})

	t.Run("invalid total to distribute, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.TotalToDistribute = "total to distribute"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "total to distribute"))
	})

	t.Run("invalid newly minted, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.TotalNewlyMinted = "total newly minted"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "total newly minted"))
	})

	t.Run("invalid rewards per block, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.RewardsPerBlock = "rewards per block"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "rewards per block"))
	})

	t.Run("invalid rewards protocol, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.RewardsForProtocolSustainability = "rewards protocol"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "rewards protocol"))
	})

	t.Run("invalid node price, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.NodePrice = "node price"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "invalid"))
		require.True(t, strings.Contains(err.Error(), "node price"))
	})

	t.Run("invalid prev epoch hash, should return error", func(t *testing.T) {
		t.Parallel()

		apiEpochStartInfoCopy := *apiEpochStartInfo
		apiEpochStartInfoCopy.PrevEpochStartHash = "prev epoch hash"

		epochStartInfo, err := esi.ProcessEpochStartInfo(&apiEpochStartInfoCopy)
		require.Nil(t, epochStartInfo)
		require.Error(t, err)
	})
}
