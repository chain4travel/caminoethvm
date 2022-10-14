package evm

import (
	"errors"
	"github.com/chain4travel/caminogo/utils/wrappers"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	percentDenominator uint64 = 1_000_000
	x2cRateUint64             = uint64(x2cRateInt64)
)

var (
	Slot0 = common.Hash{0x00}
	Slot1 = common.Hash{0x01}
)

type RewardCalculation struct {
	ValidatorRewardAmount     *big.Int `serialize:"true" json:"validatorRewardAmount"`
	ValidatorRewardToExport   uint64   `serialize:"true" json:"validatorRewardToExport"`
	IncentivePoolRewardAmount *big.Int `serialize:"true" json:"incentivePoolRewardAmount"`
	CoinbaseAmountToSub       *big.Int `serialize:"true" json:"coinbaseAmountToSub"`

	// Needed for validation that calculation can be applied
	PrevFeesBurned           *big.Int `serialize:"true" json:"prevFeesBurned"`
	PrevValidatorRewards     *big.Int `serialize:"true" json:"prevValidatorRewards"`
	PrevIncentivePoolRewards *big.Int `serialize:"true" json:"prevIncentivePoolRewards"`
}

// CalculateRewards calculates the rewards for validators and incentive pool account
//
//	feesBurned: the amount of fees burned already, balance of the coinbase address
//	validatorRewards: the amount of validator's rewards already paid out, state at slot 0 of coinbase
//	incentivePoolRewards: the amount of incentive pool's rewards already paid out, state at slot 0 of
//		incentive pool account
//	feeRewardRate: the percentage of fees to be paid out to validators, denominated in `percentDenominator`
//	incentivePoolRate: the percentage of fees to be paid out to incentive pool, denominated in `percentDenominator`
func CalculateRewards(
	feesBurned, validatorRewards, incentivePoolRewards *big.Int,
	feeRewardRate, incentivePoolRate uint64,
) (RewardCalculation, error) {
	calc := RewardCalculation{
		PrevFeesBurned:           feesBurned,
		PrevValidatorRewards:     validatorRewards,
		PrevIncentivePoolRewards: incentivePoolRewards,
	}
	errs := wrappers.Errs{}
	bigZero := big.NewInt(0)

	errs.Add(
		errIf(feeRewardRate+incentivePoolRate > percentDenominator, errors.New("feeRewardRate + incentivePoolRate > 100%")),
		errIf(feesBurned.Cmp(bigZero) < 0, errors.New("feesBurned < 0")),
		errIf(validatorRewards.Cmp(bigZero) < 0, errors.New("validatorRewards < 0")),
		errIf(incentivePoolRewards.Cmp(bigZero) < 0, errors.New("incentivePoolRewards < 0")),
	)
	if errs.Errored() {
		return calc, errs.Err
	}

	totalFeesAmount := big.NewInt(0).Add(feesBurned, incentivePoolRewards)
	totalFeesAmount.Add(totalFeesAmount, validatorRewards)

	feeRewardAmount := calculateReward(totalFeesAmount, validatorRewards, feeRewardRate)

	// Validator's reward is exported to the P-Chain, we intentionally loose precision so the C-Chain's
	// "decimal loss" can be accumulated in the account for future collections
	feeRewardAmount = bigDiv(feeRewardAmount, x2cRateUint64)
	calc.ValidatorRewardToExport = feeRewardAmount.Uint64()
	calc.ValidatorRewardAmount = bigMul(feeRewardAmount, x2cRateUint64)

	calc.IncentivePoolRewardAmount = calculateReward(totalFeesAmount, incentivePoolRewards, incentivePoolRate)

	calc.CoinbaseAmountToSub = new(big.Int).Add(calc.ValidatorRewardAmount, calc.IncentivePoolRewardAmount)

	return calc, nil
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
