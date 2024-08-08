// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/coreth/utils"
	"github.com/ethereum/go-ethereum/common"
)

// Gas Price
const (
	SunrisePhase0ExtraDataSize        = 0
	SunrisePhase0BaseFee       uint64 = 200_000_000_000
)

// Camino ChainIDs
var (
	// CaminoChainID ...
	CaminoChainID = big.NewInt(500)
	// CaminoChainID ...
	ColumbusChainID = big.NewInt(501)
	// KopernikusChainID ...
	KopernikusChainID = big.NewInt(502)
)

var (
	// CaminoChainConfig is the configuration for Camino Main Network
	CaminoChainConfig = getCaminoChainConfig(constants.CaminoID, CaminoChainID)

	// ColumbusChainConfig is the configuration for Columbus Test Network
	ColumbusChainConfig = getCaminoChainConfig(constants.ColumbusID, ColumbusChainID)

	// KopernikusChainConfig is the configuration for Kopernikus Dev Network
	KopernikusChainConfig = getCaminoChainConfig(constants.KopernikusID, KopernikusChainID)

	TestCaminoChainConfig = &ChainConfig{AvalancheContext{common.Hash{1}}, big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)}
)

func getCaminoChainConfig(networkID uint32, chainID *big.Int) *ChainConfig {
	chainConfig := getChainConfig(networkID, chainID)
	chainConfig.EIP150Hash = common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0")
	chainConfig.SunrisePhase0BlockTimestamp = getUpgradeTime(networkID, version.SunrisePhase0Times)
	chainConfig.BerlinBlockTimestamp = getUpgradeTime(networkID, version.BerlinPhaseTimes)
	return chainConfig
}

// CaminoRules returns the Camino modified rules to support Camino
// network upgrades
func (c *ChainConfig) CaminoRules(blockNum, blockTimestamp *big.Int) Rules {
	rules := c.AvalancheRules(blockNum, blockTimestamp)

	rules.IsSunrisePhase0 = c.IsSunrisePhase0(blockTimestamp)
	return rules
}

// IsSunrisePhase0 returns whether [blockTimestamp] represents a block
// with a timestamp after the Sunrise Phase 0 upgrade time.
func (c *ChainConfig) IsSunrisePhase0(time *big.Int) bool {
	return utils.IsForked(c.SunrisePhase0BlockTimestamp, time)
}

// IsBerlin returns whether [time] represents a block
// with a timestamp after the Berlin upgrade time.
func (c *ChainConfig) IsBerlin(time *big.Int) bool {
	return utils.IsForked(c.BerlinBlockTimestamp, time)
}
