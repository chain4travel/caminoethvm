package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateEverything(t *testing.T) {
	blackHoleAddressBalance := big.NewInt(1000000000000000000)
	payedOutBalance := big.NewInt(0)
	feeRewardRate := uint64(1000000)
	incentivePoolBalance := big.NewInt(0)
	incentivePoolRate := uint64(1000000)
	x2cDenominationRate := big.NewInt(1000000000000000000)

	feeRewardAmountToExport, newPayedOutBalance, newBlackHoleAddressBalance, newIncentivePoolBalance := CalculateEverything(
		*blackHoleAddressBalance,
		*payedOutBalance,
		feeRewardRate,
		*incentivePoolBalance,
		incentivePoolRate,
		*x2cDenominationRate,
	)

	assert.Equal(t, uint64(1000000000000000000), feeRewardAmountToExport)
	assert.Equal(t, big.NewInt(1000000000000000000), &newPayedOutBalance)
	assert.Equal(t, big.NewInt(0), &newBlackHoleAddressBalance)
	assert.Equal(t, big.NewInt(1000000000000000000), &newIncentivePoolBalance)
}
