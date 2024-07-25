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

	TestCaminoChainConfig = &ChainConfig{
		AvalancheContext:                AvalancheContext{common.Hash{1}},
		ChainID:                         big.NewInt(1),
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    nil,
		DAOForkSupport:                  false,
		EIP150Block:                     big.NewInt(0),
		EIP155Block:                     big.NewInt(0),
		EIP158Block:                     big.NewInt(0),
		ByzantiumBlock:                  big.NewInt(0),
		ConstantinopleBlock:             big.NewInt(0),
		PetersburgBlock:                 big.NewInt(0),
		IstanbulBlock:                   big.NewInt(0),
		MuirGlacierBlock:                big.NewInt(0),
		ApricotPhase1BlockTimestamp:     utils.NewUint64(0),
		ApricotPhase2BlockTimestamp:     utils.NewUint64(0),
		ApricotPhase3BlockTimestamp:     utils.NewUint64(0),
		ApricotPhase4BlockTimestamp:     utils.NewUint64(0),
		ApricotPhase5BlockTimestamp:     utils.NewUint64(0),
		SunrisePhase0BlockTimestamp:     utils.NewUint64(0),
		ApricotPhasePre6BlockTimestamp:  utils.NewUint64(0),
		ApricotPhase6BlockTimestamp:     utils.NewUint64(0),
		ApricotPhasePost6BlockTimestamp: utils.NewUint64(0),
		BanffBlockTimestamp:             utils.NewUint64(0),
		CortinaBlockTimestamp:           utils.NewUint64(0),
		BerlinBlockTimestamp:            utils.NewUint64(0),
		DUpgradeBlockTimestamp:          nil,
	}
)

func getCaminoChainConfig(networkID uint32, chainID *big.Int) *ChainConfig {
	chainConfig := getChainConfig(networkID, chainID)
	chainConfig.SunrisePhase0BlockTimestamp = getUpgradeTime(networkID, version.SunrisePhase0Times)
	chainConfig.BerlinBlockTimestamp = getUpgradeTime(networkID, version.BerlinPhaseTimes)
	return chainConfig
}

// CaminoRules returns the Camino modified rules to support Camino
// network upgrades
func (c *ChainConfig) CaminoRules(blockNum *big.Int, blockTimestamp uint64) Rules {
	rules := c.AvalancheRules(blockNum, blockTimestamp)

	rules.IsSunrisePhase0 = c.IsSunrisePhase0(blockTimestamp)
	return rules
}

// IsSunrisePhase0 returns whether [blockTimestamp] represents a block
// with a timestamp after the Sunrise Phase 0 upgrade time.
func (c *ChainConfig) IsSunrisePhase0(time uint64) bool {
	return utils.IsTimestampForked(c.SunrisePhase0BlockTimestamp, time)
}

// IsBerlin returns whether [time] represents a block
// with a timestamp after the Berlin upgrade time.
func (c *ChainConfig) IsBerlin(time uint64) bool {
	return utils.IsTimestampForked(c.BerlinBlockTimestamp, time)
}
