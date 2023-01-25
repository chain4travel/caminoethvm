package evm

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	gconstants "github.com/ava-labs/coreth/constants"

	"github.com/ava-labs/avalanchego/utils/units"

	"github.com/stretchr/testify/require"
)

func TestEVMStateTransfer(t *testing.T) {
	tests := []struct {
		name                   string
		balance                int64
		slotBalance            int64
		amountToDistribute     uint64
		expectedIPReward       uint64
		expectedNewBalance     uint64
		expectedNewSlotBalance uint64
	}{
		{
			name:                   "Happy path: coinbase 10, slot 0 => export 3, incentive 3, coinbase 4, slot 4",
			balance:                10 * int64(units.Avax),
			slotBalance:            0,
			amountToDistribute:     10,
			expectedIPReward:       3 * units.Avax,
			expectedNewBalance:     4 * units.Avax,
			expectedNewSlotBalance: 4 * units.Avax,
		},
		{
			name:                   "Happy path: coinbase 11, slot 0 => export 3, incentive 3, coinbase 5, slot 4",
			balance:                11 * int64(units.Avax),
			slotBalance:            0,
			amountToDistribute:     11,
			expectedIPReward:       3 * units.Avax,
			expectedNewBalance:     5 * units.Avax,
			expectedNewSlotBalance: 4 * units.Avax,
		},
		{
			name:                   "Happy path: bigger numbers, ratio without decimal loss, BH balance == Balance slot",
			balance:                100 * int64(units.Avax),
			slotBalance:            10 * int64(units.Avax),
			amountToDistribute:     90,
			expectedIPReward:       27 * units.Avax,
			expectedNewBalance:     46 * units.Avax, // == 90 - 2 * 27 + 10
			expectedNewSlotBalance: 46 * units.Avax, // == 40% of 90 + 10
		},
		{
			name:                   "Happy path: bigger numbers, ratio with decimal loss, BH balance > Balance slot",
			balance:                90 * int64(units.Avax),
			slotBalance:            6 * int64(units.Avax),
			amountToDistribute:     84,
			expectedIPReward:       25 * units.Avax,
			expectedNewBalance:     40 * units.Avax, // == 84 - 2 * 25 + 6
			expectedNewSlotBalance: 39 * units.Avax, // == 40% of 84 + 6
		},
		{
			name:                   "Simulation block 1: fee increase: 100, coinbase (10)0, slot 0",
			balance:                100 * int64(units.Avax),
			slotBalance:            0 * int64(units.Avax),
			amountToDistribute:     100,
			expectedIPReward:       30 * units.Avax,
			expectedNewBalance:     40 * units.Avax,
			expectedNewSlotBalance: 40 * units.Avax,
		},
		{
			name:                   "Simulation block 2: fee increase: 90, coinbase 130, slot 40",
			balance:                130 * int64(units.Avax),
			slotBalance:            40 * int64(units.Avax),
			amountToDistribute:     90,
			expectedIPReward:       27 * units.Avax,
			expectedNewBalance:     76 * units.Avax,
			expectedNewSlotBalance: 76 * units.Avax,
		},
		{
			name:                   "Simulation block 3: fee increase: 84, coinbase 160, slot 76",
			balance:                160 * int64(units.Avax),
			slotBalance:            76 * int64(units.Avax),
			amountToDistribute:     84,
			expectedIPReward:       25 * units.Avax,
			expectedNewBalance:     110 * units.Avax,
			expectedNewSlotBalance: 109 * units.Avax,
		},
		{
			name:                   "Simulation block 4: fee increase: 85, coinbase 196, slot 76",
			balance:                195 * int64(units.Avax),
			slotBalance:            109 * int64(units.Avax),
			amountToDistribute:     86,
			expectedIPReward:       25 * units.Avax,
			expectedNewBalance:     145 * units.Avax,
			expectedNewSlotBalance: 142 * units.Avax,
		},
		{
			name:                   "Simulation block 5: fee increase: 84, coinbase 196, slot 76",
			balance:                229 * int64(units.Avax),
			slotBalance:            142 * int64(units.Avax),
			amountToDistribute:     87,
			expectedIPReward:       26 * units.Avax,
			expectedNewBalance:     177 * units.Avax,
			expectedNewSlotBalance: 176 * units.Avax,
		},
		{
			name:                   "Simulation block 1-5 in one go: fee increase: 100+90+84+85+84, coinbase 443, slot 0",
			balance:                443 * int64(units.Avax),
			slotBalance:            0 * int64(units.Avax),
			amountToDistribute:     443,
			expectedIPReward:       132 * units.Avax, // -1 to the accumulated rewards (133) from the previous 👆 blocks
			expectedNewBalance:     179 * units.Avax, // which most likely will be aligned, because of +2 in BH balance
			expectedNewSlotBalance: 176 * units.Avax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, vm, _, _, _ := GenesisVM(t, true, genesisJSONSunrisePhase0, "", "")
			state, err := vm.blockChain.State()
			require.NoError(t, err)

			// Add balance to coinbase address
			state.AddBalance(gconstants.BlackholeAddr, big.NewInt(tt.balance))

			// Add slot balance to coinbase address
			state.SetState(gconstants.BlackholeAddr, BalanceSlot, common.BigToHash(big.NewInt(tt.slotBalance)))

			// Cal the rewards tx
			tx, err := vm.NewCollectRewardsTx(tt.amountToDistribute, 1)
			require.NoError(t, err)

			err = tx.EVMStateTransfer(vm.ctx, state)
			require.NoError(t, err)

			// assert incentive balance
			incentiveBalance := state.GetBalance(common.Address(FeeRewardAddressID)).Uint64()
			require.Equal(t, tt.expectedIPReward, incentiveBalance, fmt.Sprintf("expected %d, got (actual) %d", tt.expectedIPReward, incentiveBalance))

			// assert coinbase balance
			newCoinbaseBalance := state.GetBalance(gconstants.BlackholeAddr).Uint64()
			require.Equal(t, tt.expectedNewBalance, newCoinbaseBalance, fmt.Sprintf("expected %d, got (actual) %d", tt.expectedNewBalance, newCoinbaseBalance))

			// assert slot balance
			newSlotBalance := state.GetState(gconstants.BlackholeAddr, BalanceSlot).Big().Uint64()
			require.Equal(t, tt.expectedNewSlotBalance, newSlotBalance, fmt.Sprintf("expected %d, got (actual) %d", tt.expectedNewSlotBalance, newSlotBalance))
		})
	}
}
