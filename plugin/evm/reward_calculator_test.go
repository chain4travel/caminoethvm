package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewardRatesExceed100Percent(t *testing.T) {
	blackHoleAddressBalance := big.NewInt(1_000_000_000_000_000_000)
	payedOutBalance := big.NewInt(0)
	incentivePoolBalance := big.NewInt(0)
	feeRewardRate := uint64(0.51 * float64(percentDenominator))
	incentivePoolRate := uint64(0.50 * float64(percentDenominator))

	_, err := CalculateRewards(
		blackHoleAddressBalance,
		payedOutBalance,
		incentivePoolBalance,
		feeRewardRate,
		incentivePoolRate,
	)

	assert.Error(t, err)
	assert.Equal(t, "feeRewardRate + incentivePoolRate > 100%", err.Error())
}

func TestCalculate10PercentReward(t *testing.T) {
	blackHoleAddressBalance := big.NewInt(1_000_000_000_000_000_000)
	payedOutBalance := big.NewInt(0)
	incentivePoolBalance := big.NewInt(0)
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	calc, err := CalculateRewards(
		blackHoleAddressBalance,
		payedOutBalance,
		incentivePoolBalance,
		feeRewardRate,
		incentivePoolRate,
	)

	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(100_000_000_000_000_000), calc.ValidatorRewardAmount)
	assert.Equal(t, uint64(100_000_000), calc.ValidatorRewardToExport)
	assert.Equal(t, big.NewInt(100_000_000_000_000_000), calc.IncentivePoolRewardAmount)
	assert.Equal(t, big.NewInt(200_000_000_000_000_000), calc.CoinbaseAmountToSub)
}
