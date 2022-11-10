// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"fmt"

	"github.com/ava-labs/coreth/core/state"
	"github.com/ethereum/go-ethereum/log"
)

// calculateAndCollectRewards calculates the rewards and issues the CollectRewardsTx
// Errors are logged and ignored as they are not affecting the other actions
func (b *Block) calculateAndCollectRewards() {
	state, err := b.vm.blockChain.State()
	if err != nil {
		return
	}
	calc, err := b.calculateRewards(state)
	if err != nil {
		return
	}

	tx, err := b.createRewardsCollectionTx(calc)
	if err != nil {
		log.Info("Issuing of the rewards collection skipped", "error", err)
	} else {
		log.Info("Issuing of the rewards collection tx", "txID", tx.ID(), "rewards to export", calc.ValidatorRewardToExport)
		b.vm.issueTx(tx, true /*=local*/)
	}
}

func (b *Block) calculateRewards(state *state.StateDB) (*RewardCalculation, error) {
	header := b.ethBlock.Header()

	feesBurned := state.GetBalance(header.Coinbase)
	validatorRewards := state.GetState(header.Coinbase, Slot1).Big()
	incentivePoolRewards := state.GetState(header.Coinbase, Slot2).Big()

	calculation, err := CalculateRewards(
		feesBurned,
		validatorRewards,
		incentivePoolRewards,
		FeeRewardRate,
		IncentivePoolRewardRate,
	)

	if err != nil {
		return nil, err
	}

	if calculation.ValidatorRewardToExport < FeeRewardMinAmountToExport {
		return nil, fmt.Errorf("calculated fee reward amount %d is less than the minimum amount to export", calculation.ValidatorRewardToExport)
	}

	return &calculation, nil
}

func (b *Block) createRewardsCollectionTx(calculation *RewardCalculation) (*Tx, error) {
	if calculation == nil {
		return nil, fmt.Errorf("no rewards to collect")
	}

	h := b.ethBlock.Header()
	tx, err := b.vm.NewCollectRewardsTx(
		calculation.Result(),
		b.ID(),
		b.ethBlock.Time(),
		h.Coinbase,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
