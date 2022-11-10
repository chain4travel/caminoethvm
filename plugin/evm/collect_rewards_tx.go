// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/verify"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	FeeRewardMinAmountToExport    = uint64(200_000)
	FeeRewardExportAddressStr     = "0x010000000000000000000000000000000000000c"
	IncentivePoolRewardAddressStr = "0x010000000000000000000000000000000000000c"

	// Assumption: Following rates are denominated to uint64 from floating point ratio (ratio * evm.percentDenominator)
	FeeRewardRate           = uint64(300_000)
	IncentivePoolRewardRate = uint64(300_000)
)

var (
	_                           UnsignedAtomicTx = &UnsignedCollectRewardsTx{}
	FeeRewardExportAddress                       = common.HexToAddress(FeeRewardExportAddressStr)
	IncentivePoolRewardAddress                   = common.HexToAddress(IncentivePoolRewardAddressStr)
	FeeRewardExportAddressId, _                  = ids.ToShortID(FeeRewardExportAddress.Bytes())
)

// UnsignedCollectRewardsTx is an internal rewards collection transaction
type UnsignedCollectRewardsTx struct {
	avax.Metadata
	// ID of the network on which this tx was issued
	NetworkID uint32 `serialize:"true" json:"networkID"`
	// ID of this blockchain.
	BlockchainID ids.ID `serialize:"true" json:"blockchainID"`
	// Which chain to send the funds to
	DestinationChain ids.ID `serialize:"true" json:"destinationChain"`
	// Outputs that are exported to the chain
	ExportedOutputs []*avax.TransferableOutput `serialize:"true" json:"exportedOutputs"`

	BlockId        ids.ID `serialize:"true" json:"blockId"`
	BlockTimestamp uint64 `serialize:"true" json:"blockTimestamp"`

	RewardCalculation RewardCalculationResult `serialize:"true" json:"rewardCalculation"`

	Coinbase common.Address `serialize:"true" json:"coinbase"`
}

// InputUTXOs returns a set of all the hash(address:nonce) exporting funds.
func (tx *UnsignedCollectRewardsTx) InputUTXOs() set.Set[ids.ID] {
	// Not sure if it will be needed - mock
	return set.NewSet[ids.ID](0)
}

// Verify this transaction is well-formed
func (tx *UnsignedCollectRewardsTx) Verify(
	ctx *snow.Context,
	_rules params.Rules,
) error {
	switch {
	case tx == nil:
		return errNilTx
	case tx.NetworkID != ctx.NetworkID:
		return errWrongNetworkID
	case ctx.ChainID != tx.BlockchainID:
		return errWrongBlockchainID
	}

	if err := tx.RewardCalculation.Verify(); err != nil {
		return err
	}

	// Make sure that the tx has a valid peer chain ID
	if err := verify.SameSubnet(ctx, tx.DestinationChain); err != nil {
		return errWrongChainID
	}
	if tx.DestinationChain != constants.PlatformChainID {
		return errWrongChainID
	}

	return nil
}

func (tx *UnsignedCollectRewardsTx) GasUsed(bool) (uint64, error) {
	return 0, nil
}

// Amount of [assetID] burned by this transaction
func (tx *UnsignedCollectRewardsTx) Burned(_assetID ids.ID) (uint64, error) {
	// Let me lie here
	return 0, nil
}

// SemanticVerify this transaction is valid.
func (tx *UnsignedCollectRewardsTx) SemanticVerify(
	vm *VM,
	_stx *Tx,
	b *Block,
	_baseFee *big.Int,
	rules params.Rules,
) error {
	if err := tx.Verify(vm.ctx, rules); err != nil {
		return err
	}

	stateDB, err := vm.blockChain.State()
	if err != nil {
		log.Error("Cannot access current EVM stateDB", "error", err)
		return fmt.Errorf("cannot access current EVM stateDB: %w", err)
	}
	if tx.Coinbase != b.ethBlock.Coinbase() {
		return fmt.Errorf("coinbase address mismatch Tx %s vs Block's %s", tx.Coinbase.Hex(), b.ethBlock.Coinbase().Hex())
	}

	calculation := tx.RewardCalculation.Calculation()

	currValidatorRewardPayedOut := stateDB.GetState(tx.Coinbase, Slot1).Big()
	if calculation.PrevValidatorRewards.Cmp(currValidatorRewardPayedOut) != 0 {
		log.Info("validator rewards mismatch", "prevFeesBurned", calculation.PrevFeesBurned, "currValidatorRewardPayedOut", currValidatorRewardPayedOut)
		return fmt.Errorf("validator rewards mismatch")
	}

	ipRewardsPayedOut := stateDB.GetState(tx.Coinbase, Slot2).Big()
	if calculation.PrevIncentivePoolRewards.Cmp(ipRewardsPayedOut) != 0 {
		log.Info("Incentive pool rewards mismatch", "prevIncentivePoolRewards", calculation.PrevIncentivePoolRewards, "ipRewardsPayedOut", ipRewardsPayedOut)
		return fmt.Errorf("incentive pool balance mismatch")
	}

	// 1. Redo the calculation and compare the results
	prevCalc, err := CalculateRewards(calculation.PrevFeesBurned, calculation.PrevValidatorRewards, calculation.PrevIncentivePoolRewards, FeeRewardRate, IncentivePoolRewardRate)
	if err != nil {
		return fmt.Errorf("cannot repeat the reward calculation on previous state: %w", err)
	}

	if calculation.ValidatorRewardAmount.Cmp(prevCalc.ValidatorRewardAmount) != 0 ||
		calculation.IncentivePoolRewardAmount.Cmp(prevCalc.IncentivePoolRewardAmount) != 0 ||
		calculation.ValidatorRewardToExport != prevCalc.ValidatorRewardToExport ||
		calculation.CoinbaseAmountToSub.Cmp(prevCalc.CoinbaseAmountToSub) != 0 {
		return fmt.Errorf("repeated reward calculation on previous state does not match the Tx")
	}

	// 2. Check that the state is correct - do the calculation on current state and check calculated rewards are not less than in the Tx
	currFeesBurned := stateDB.GetBalance(tx.Coinbase)
	if currFeesBurned.Cmp(calculation.PrevFeesBurned) < 0 {
		return fmt.Errorf("current Coinbase balance is less than the balance the CollectRewardsTx was issued for")
	}

	currCalc, err := CalculateRewards(currFeesBurned, currValidatorRewardPayedOut, ipRewardsPayedOut, FeeRewardRate, IncentivePoolRewardRate)
	if err != nil {
		return fmt.Errorf("cannot repeat the reward calculation on current state: %w", err)
	}

	if calculation.ValidatorRewardAmount.Cmp(currCalc.ValidatorRewardAmount) > 0 ||
		calculation.IncentivePoolRewardAmount.Cmp(currCalc.IncentivePoolRewardAmount) > 0 ||
		calculation.ValidatorRewardToExport > currCalc.ValidatorRewardToExport ||
		calculation.CoinbaseAmountToSub.Cmp(currCalc.CoinbaseAmountToSub) > 0 {
		return fmt.Errorf("repeated reward calculation on current state does not follow the requirements")
	}

	// 3. Check the UTXO's owner, amount & currency
	if len(tx.ExportedOutputs) != 1 {
		return fmt.Errorf("expected single exported output, got %d", len(tx.ExportedOutputs))
	}

	eo := tx.ExportedOutputs[0]
	if eo.AssetID() != vm.ctx.AVAXAssetID {
		return fmt.Errorf("expected AVAX asset in ExportedOutput, got %s", eo.AssetID())
	}

	if eo.Out.Amount() != calculation.ValidatorRewardToExport {
		return fmt.Errorf("expected %d AVAX in ExportedOutput, got %d", calculation.ValidatorRewardToExport, eo.Out.Amount())
	}

	out, ok := eo.Out.(*secp256k1fx.TransferOutput)
	if !ok {
		return fmt.Errorf("expected secp256k1fx.TransferOutput in ExportedOutput, got %T", eo.Out)
	}

	if len(out.OutputOwners.Addrs) != 1 {
		return fmt.Errorf("expected single output owner in ExportedOutput")
	}
	if err != nil {
		return fmt.Errorf("failed to get short ID for fee reward export address: %w", err)
	}
	if out.OutputOwners.Addrs[0] != FeeRewardExportAddressId {
		return fmt.Errorf("expected %s as output owner of the fee reward, got %s", FeeRewardExportAddressStr, out.OutputOwners.Addrs[0].Hex())
	}

	return nil
}

// AtomicOps returns the atomic operations for this transaction.
func (tx *UnsignedCollectRewardsTx) AtomicOps() (ids.ID, *atomic.Requests, error) {
	txID := tx.ID()

	out := tx.ExportedOutputs[0]
	utxo := &avax.UTXO{
		UTXOID: avax.UTXOID{
			TxID:        txID,
			OutputIndex: uint32(0),
		},
		Asset: avax.Asset{ID: out.AssetID()},
		Out:   out.Out,
	}

	utxoBytes, err := Codec.Marshal(codecVersion, utxo)
	if err != nil {
		return ids.ID{}, nil, err
	}
	utxoID := utxo.InputID()
	elem := &atomic.Element{
		Key:   utxoID[:],
		Value: utxoBytes,
	}
	if out, ok := utxo.Out.(avax.Addressable); ok {
		elem.Traits = out.Addresses()
	}

	elems := []*atomic.Element{elem}
	return tx.DestinationChain, &atomic.Requests{PutRequests: elems}, nil
}

func (vm *VM) NewCollectRewardsTx(
	calculation RewardCalculationResult,
	blockId ids.ID,
	blockTimestamp uint64,
	coinbase common.Address,
) (*Tx, error) {
	// Create the transaction
	utx := &UnsignedCollectRewardsTx{
		NetworkID:         vm.ctx.NetworkID,
		BlockchainID:      vm.ctx.ChainID,
		DestinationChain:  constants.PlatformChainID,
		BlockId:           blockId,
		BlockTimestamp:    blockTimestamp,
		RewardCalculation: calculation,
		Coinbase:          coinbase,
	}

	utx.ExportedOutputs = []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: vm.ctx.AVAXAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: calculation.ValidatorRewardToExport,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  0,
				Threshold: 1,
				Addrs:     []ids.ShortID{FeeRewardExportAddressId},
			},
		},
	}}

	tx := &Tx{UnsignedAtomicTx: utx}
	if err := tx.Sign(vm.codec, nil); err != nil {
		return nil, err
	}

	return tx, utx.Verify(vm.ctx, vm.currentRules())
}

// EVMStateTransfer executes the state update from the atomic export transaction
func (tx *UnsignedCollectRewardsTx) EVMStateTransfer(ctx *snow.Context, state *state.StateDB) error {
	calculation := tx.RewardCalculation.Calculation()
	state.SubBalance(tx.Coinbase, calculation.CoinbaseAmountToSub)
	validatorRewards := new(big.Int).Add(calculation.PrevValidatorRewards, calculation.ValidatorRewardAmount)

	ipRewards := new(big.Int).Add(calculation.PrevIncentivePoolRewards, calculation.IncentivePoolRewardAmount)
	state.AddBalance(IncentivePoolRewardAddress, calculation.IncentivePoolRewardAmount)

	state.SetState(tx.Coinbase, Slot1, common.BigToHash(validatorRewards))
	state.SetState(tx.Coinbase, Slot2, common.BigToHash(ipRewards))

	return nil
}
