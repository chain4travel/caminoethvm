package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ZERO = big.NewInt(0)
	CAM1 = big.NewInt(1_000_000_000_000_000_000)
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

func TestCalculationWithoutFeeIncrementIsZeroRewards(t *testing.T) {
	// sytuacja: brak nowych feesBurned
	// oczekiwany wynik: brak nagród
	feesBurned := CAM1
	validatorRewards := ZERO
	incentivePoolRewards := ZERO
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	calc1, err := CalculateRewards(
		feesBurned,
		validatorRewards,     // 10% feesów ale z denominacją
		incentivePoolRewards, //10% feesów
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	coinbaseAmountAfterCalc1 := big.NewInt(0).Sub(feesBurned, calc1.CoinbaseAmountToSub)
	validatorAmountAfterCalc1 := big.NewInt(0).Add(validatorRewards, calc1.ValidatorRewardAmount)
	incentivePoolRewardAmountAfterCalc1 := big.NewInt(0).Add(incentivePoolRewards, calc1.IncentivePoolRewardAmount)

	calc2, err := CalculateRewards(
		coinbaseAmountAfterCalc1,
		validatorAmountAfterCalc1,
		incentivePoolRewardAmountAfterCalc1,
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	assert.Equal(t, ZERO, calc2.CoinbaseAmountToSub)
	assert.Equal(t, ZERO, calc2.ValidatorRewardAmount)
	assert.Equal(t, ZERO, calc2.IncentivePoolRewardAmount)
}

func TestSumOfCalculationsEqualsTotalCalculation(t *testing.T) {
	feesBurned := CAM1
	validatorRewards := ZERO
	incentivePoolRewards := ZERO
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	calc1, err := CalculateRewards(
		feesBurned,
		validatorRewards,     // 10% feesów ale z denominacją
		incentivePoolRewards, //10% feesów
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	// po każdym bloku dodajmy 1 CAM (1_000_000_000_000_000_000)

	coinbaseAmountAfterCalc1 := big.NewInt(0).Sub(feesBurned, calc1.CoinbaseAmountToSub)
	validatorAmountAfterCalc1 := big.NewInt(0).Add(validatorRewards, calc1.ValidatorRewardAmount)
	incentivePoolRewardAmountAfterCalc1 := big.NewInt(0).Add(incentivePoolRewards, calc1.IncentivePoolRewardAmount)
	// niezmiennik: po kazdej kalkulacji calc.CoinbaseAmountToSub ==  calc.ValidatorRewardAmount + calc.IncentivePoolRewardAmount
	// total = coinbaseAmountAfterCalc + validatorAmountAfterCalc + incentivePoolAmountAfterCalc

	total := CAM1
	assert.Equal(t,
		total,
		bigAdd(coinbaseAmountAfterCalc1, validatorAmountAfterCalc1, incentivePoolRewardAmountAfterCalc1),
	)
	feesBurned = bigAdd(coinbaseAmountAfterCalc1, CAM1)
	calc2, err := CalculateRewards(
		feesBurned,
		validatorAmountAfterCalc1,
		incentivePoolRewardAmountAfterCalc1,
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	coinbaseAmountAfterCalc2 := big.NewInt(0).Sub(feesBurned, calc2.CoinbaseAmountToSub)
	validatorAmountAfterCalc2 := bigAdd(validatorAmountAfterCalc1, calc2.ValidatorRewardAmount)
	incentivePoolRewardAmountAfterCalc2 := bigAdd(incentivePoolRewardAmountAfterCalc1, calc2.IncentivePoolRewardAmount)

	total = bigMul(CAM1, 2)
	assert.Equal(t,
		total,
		bigAdd(coinbaseAmountAfterCalc2, validatorAmountAfterCalc2, incentivePoolRewardAmountAfterCalc2),
	)

	feesBurned = bigAdd(coinbaseAmountAfterCalc2, CAM1)
	calc3, err := CalculateRewards(
		feesBurned,
		validatorAmountAfterCalc2,
		incentivePoolRewardAmountAfterCalc2,
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	coinbaseAmountAfterCalc3 := big.NewInt(0).Sub(feesBurned, calc3.CoinbaseAmountToSub)
	validatorAmountAfterCalc3 := bigAdd(validatorAmountAfterCalc2, calc3.ValidatorRewardAmount)
	incentivePoolRewardAmountAfterCalc3 := bigAdd(incentivePoolRewardAmountAfterCalc2, calc3.IncentivePoolRewardAmount)

	total = bigMul(CAM1, 3)
	assert.Equal(t,
		total,
		bigAdd(coinbaseAmountAfterCalc3, validatorAmountAfterCalc3, incentivePoolRewardAmountAfterCalc3),
	)

	calcTotal, err := CalculateRewards(
		bigMul(CAM1, 3),
		ZERO,
		ZERO,
		feeRewardRate,
		incentivePoolRate,
	)
	assert.NoError(t, err)

	assert.Equal(t, incentivePoolRewardAmountAfterCalc3, validatorAmountAfterCalc3)
	assert.Equal(t, calcTotal.ValidatorRewardAmount, validatorAmountAfterCalc3)
	assert.Equal(t, calcTotal.IncentivePoolRewardAmount, incentivePoolRewardAmountAfterCalc3)
}

func TestNegativeBlackHoleAddressBalance(t *testing.T) {
	feesBurned := big.NewInt(-1)
	validatorRewards := big.NewInt(0)
	incentivePoolRewards := big.NewInt(0)
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	_, err := CalculateRewards(
		feesBurned,
		validatorRewards,
		incentivePoolRewards,
		feeRewardRate,
		incentivePoolRate,
	)

	assert.Error(t, err)
	assert.Equal(t, "feesBurned < 0", err.Error())
}

func TestNegativePayedOutBalance(t *testing.T) {
	feesBurned := big.NewInt(1_000_000_000_000_000_000)
	validatorRewards := big.NewInt(-1)
	incentivePoolRewards := big.NewInt(0)
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	_, err := CalculateRewards(
		feesBurned,
		validatorRewards,
		incentivePoolRewards,
		feeRewardRate,
		incentivePoolRate,
	)

	assert.Error(t, err)
	assert.Equal(t, "validatorRewards < 0", err.Error())
}

func TestNegativeIncentivePoolBalance(t *testing.T) {
	feesBurned := big.NewInt(1_000_000_000_000_000_000)
	validatorRewards := big.NewInt(0)
	incentivePoolRewards := big.NewInt(-1)
	feeRewardRate := uint64(0.10 * float64(percentDenominator))
	incentivePoolRate := uint64(0.10 * float64(percentDenominator))

	_, err := CalculateRewards(
		feesBurned,
		validatorRewards,
		incentivePoolRewards,
		feeRewardRate,
		incentivePoolRate,
	)

	assert.Error(t, err)
	assert.Equal(t, "incentivePoolRewards < 0", err.Error())
}

func bigAdd(ingredients ...*big.Int) *big.Int {
	sum := big.NewInt(0)
	for _, ingredient := range ingredients {
		sum.Add(sum, ingredient)
	}
	return sum
}
