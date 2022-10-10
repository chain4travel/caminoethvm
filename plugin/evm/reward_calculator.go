package evm

import (
	"errors"
	"github.com/chain4travel/caminogo/utils/wrappers"
	"math/big"
)

const (
	percentDenominator uint64 = 1_000_000
	x2cRateUint64             = uint64(x2cRateInt64)
)

type RewardCalculation struct {
	FeeRewardAmountToExport       uint64
	IncentivePoolRewardToTransfer *big.Int
	NewPayedOutBalance            *big.Int
	NewBlackHoleAddressBalance    *big.Int
	NewIncentivePoolBalance       *big.Int
}

func CalculateRewards(
	blackHoleAddressBalance, payedOutBalance, incentivePoolBalance *big.Int,
	feeRewardRate, incentivePoolRate uint64,
) (RewardCalculation, error) {
	calculations := RewardCalculation{}
	errs := wrappers.Errs{}
	bigZero := big.NewInt(0)

	errs.Add(
		errIf(feeRewardRate+incentivePoolRate > percentDenominator, errors.New("feeRewardRate + incentivePoolRate > 100%")),
		errIf(blackHoleAddressBalance.Cmp(bigZero) < 0, errors.New("blackHoleAddressBalance < 0")),
		errIf(payedOutBalance.Cmp(bigZero) < 0, errors.New("payedOutBalance < 0")),
		errIf(incentivePoolBalance.Cmp(bigZero) < 0, errors.New("incentivePoolBalance < 0")),
	)
	if errs.Errored() {
		return calculations, errs.Err
	}

	// 1.1 Calculate totalFeeAmount = BHBalance + IPBalance + payedOutBalance
	totalBalance := big.NewInt(0).Add(blackHoleAddressBalance, incentivePoolBalance)
	totalBalance.Add(totalBalance, payedOutBalance)

	// 1.2 Calculate the validator reward partFeeRewardAmount = feeRewardRatio * totalFeeAmount - payedOutBalance
	feeRewardAmount := calculateReward(totalBalance, payedOutBalance, feeRewardRate)

	// 1.3 Denominate it from C-chain to P-chain precision feeRewardToExport = denominateCtoP(feeRewardAmount)
	calculations.FeeRewardAmountToExport = bigDiv(feeRewardAmount, x2cRateUint64).Uint64()

	// 1.4 Calculate payedOut = denominatePtoC(feeRewardToExport). Note: intentionally loosing precision here
	payedOut := bigMul(
		big.NewInt(int64(calculations.FeeRewardAmountToExport)),
		x2cRateUint64)

	// 1.5 Increase payedOutBalance += payedOut
	calculations.NewPayedOutBalance = new(big.Int).Add(payedOutBalance, payedOut)

	// 2.1 Calculate the Incentive Pool partIncentivePoolAmount = incentivePoolRatio * totalFeeAmount - incentivePoolBalance
	calculations.IncentivePoolRewardToTransfer = calculateReward(totalBalance, incentivePoolBalance, incentivePoolRate)

	// 2.2 Increase incentivePoolBalance += incentivePoolReward
	calculations.NewIncentivePoolBalance = new(big.Int).Add(incentivePoolBalance, calculations.IncentivePoolRewardToTransfer)

	// 3.1 Decrease BHBalance -= (payedOut + ipAmount)
	calculations.NewBlackHoleAddressBalance = new(big.Int).Sub(
		blackHoleAddressBalance,
		new(big.Int).Add(payedOut, calculations.IncentivePoolRewardToTransfer),
	)

	return calculations, nil
}

func calculateReward(total, alreadyPayed *big.Int, percentRate uint64) *big.Int {
	reward := new(big.Int).Sub(
		bigMul(total, percentRate),
		bigMul(alreadyPayed, percentDenominator),
	)

	// denominate down from percentage
	return bigDiv(reward, percentDenominator)
}

// bigMul: multiply value by rate
func bigMul(value *big.Int, rate uint64) *big.Int {
	return new(big.Int).Mul(value, big.NewInt(int64(rate)))
}

// bigDiv: divide value by rate
func bigDiv(value *big.Int, rate uint64) *big.Int {
	return new(big.Int).Div(value, big.NewInt(int64(rate)))
}

func errIf(cond bool, err error) error {
	if cond {
		return err
	}
	return nil
}
