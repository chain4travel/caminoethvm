// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/coreth/utils"

	"github.com/ethereum/go-ethereum/common"
)

// Gas Price
const (
	SunrisePhase0ExtraDataSize        = 0
	SunrisePhase0BaseFee       uint64 = 200_000_000_000
)

var (
	// CaminoChainConfig is the configuration for Camino Main Network
	CaminoChainConfig = &ChainConfig{
		ChainID:                         CaminoChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
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
		// TODO Add Cortina timestamps
	}

	// ColumbusChainConfig is the configuration for Columbus Test Network
	ColumbusChainConfig = &ChainConfig{
		ChainID:                         ColumbusChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
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
		// TODO Add Cortina timestamps
	}

	// KopernikusChainConfig is the configuration for Kopernikus Dev Network
	KopernikusChainConfig = &ChainConfig{
		ChainID:                         KopernikusChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
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
		// TODO Add Cortina timestamps
	}
)

// CaminoRules returns the Camino modified rules to support Camino
// network upgrades
func (c *ChainConfig) CaminoRules(blockNum *big.Int, blockTimestamp uint64) Rules {
	rules := c.AvalancheRules(blockNum, blockTimestamp)

	rules.IsSunrisePhase0 = c.IsSunrisePhase0(blockTimestamp)
	return rules
}
