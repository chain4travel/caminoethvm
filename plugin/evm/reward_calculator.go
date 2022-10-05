package evm

import "math/big"

const (
	percentDenominator int64 = 1_000_000
)

func CalculateEverything(
	blackHoleAddressBalance big.Int,
	payedOutBalance big.Int,
	feeRewardRate uint64,
	incentivePoolBalance big.Int,
	incentivePoolRate uint64,
	x2cDenominationRate big.Int,
) (
	feeRewardAmountToExport uint64,
	newPayedOutBalance big.Int,
	newBlackHoleAddressBalance big.Int,
	newIncentivePoolBalance big.Int,
) {
	percentDenom := big.NewInt(percentDenominator)
	// 1.1 Calculate totalFeeAmount = BHBalance + IPBalance + payedOutBalance
	totalBalance := big.NewInt(0).Add(&blackHoleAddressBalance, &incentivePoolBalance)
	totalBalance.Add(totalBalance, &payedOutBalance)

	// 1.2 Calculate the validator reward partFeeRewardAmount = feeRewardRatio * totalFeeAmount - payedOutBalance
	partFeeRewardAmount := new(big.Int).Mul(totalBalance, big.NewInt(int64(feeRewardRate)))
	feeRewardAmount := new(big.Int).Sub(
		partFeeRewardAmount,
		new(big.Int).Mul(&payedOutBalance, percentDenom),
	)

	// 1.3 Denominate it from C-chain to P-chain precision feeRewardToExport = denominateCtoP(feeRewardAmount)
	feeRewardAmount.Div(feeRewardAmount, percentDenom)
	feeRewardAmountToExport = new(big.Int).Div(feeRewardAmount, &x2cDenominationRate).Uint64()

	// 1.4 Calculate payedOut = denominatePtoC(feeRewardToExport). Note: intentionally loosing precision here
	payedOut := *new(big.Int).Mul(
		big.NewInt(int64(feeRewardAmountToExport)),
		&x2cDenominationRate)

	// 1.5 Increase payedOutBalance += payedOut
	newPayedOutBalance = *new(big.Int).Add(&payedOutBalance, &payedOut)

	// 2.1 Calculate the Incentive Pool partIncentivePoolAmount = incentivePoolRatio * totalFeeAmount - incentivePoolBalance
	partIncentivePoolAmount := new(big.Int).Mul(totalBalance, big.NewInt(int64(incentivePoolRate)))
	incentivePoolAmount := new(big.Int).Sub(
		partIncentivePoolAmount,
		new(big.Int).Mul(&incentivePoolBalance, percentDenom),
	)

	// 2.2 Increase incentivePoolBalance += incentivePoolAmount
	incentivePoolAmount.Div(incentivePoolAmount, percentDenom)
	newIncentivePoolBalance = *new(big.Int).Add(&incentivePoolBalance, incentivePoolAmount)

	// 3.1 Decrease BHBalance -= (payedOut + ipAmount)
	newBlackHoleAddressBalance = *new(big.Int).Sub(
		&blackHoleAddressBalance,
		new(big.Int).Add(&payedOut, incentivePoolAmount),
	)

	return
}
