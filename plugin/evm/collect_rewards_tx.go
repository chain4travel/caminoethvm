// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
//
// This file is a derived work, based on ava-labs code whose
// original notices appear below.
//
// It is distributed under the same license conditions as the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********************************************************

// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"github.com/chain4travel/caminogo/vms/components/verify"
	"github.com/chain4travel/caminogo/vms/secp256k1fx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"

	"github.com/chain4travel/caminoethvm/core/state"
	"github.com/chain4travel/caminoethvm/params"

	"github.com/chain4travel/caminogo/chains/atomic"
	"github.com/chain4travel/caminogo/ids"
	"github.com/chain4travel/caminogo/snow"
	"github.com/chain4travel/caminogo/utils/constants"
	"github.com/chain4travel/caminogo/vms/components/avax"
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

	RewardCalculation RewardCalculation `serialize:"true" json:"rewardCalculation"`

	Coinbase                   common.Address `serialize:"true" json:"coinbase"`
	ValidatorRewardAddress     common.Address `serialize:"true" json:"validatorRewardAddress"`
	IncentivePoolRewardAddress common.Address `serialize:"true" json:"incentivePoolRewardAddress"`
}

// InputUTXOs returns a set of all the hash(address:nonce) exporting funds.
func (tx *UnsignedCollectRewardsTx) InputUTXOs() ids.Set {
	// Not sure if it will be needed - mock
	return ids.NewSet(0)
}

// Verify this transaction is well-formed
func (tx *UnsignedCollectRewardsTx) Verify(
	ctx *snow.Context,
	rules params.Rules,
) error {
	switch {
	case tx == nil:
		return errNilTx
	case tx.NetworkID != ctx.NetworkID:
		return errWrongNetworkID
	case ctx.ChainID != tx.BlockchainID:
		return errWrongBlockchainID
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
func (tx *UnsignedCollectRewardsTx) Burned(assetID ids.ID) (uint64, error) {
	// Not sure it will be needed - mock
	return 0, nil
}

// SemanticVerify this transaction is valid.
func (tx *UnsignedCollectRewardsTx) SemanticVerify(
	vm *VM,
	stx *Tx,
	_ *Block,
	baseFee *big.Int,
	rules params.Rules,
) error {

	// TODO: Verify the fee calculation results here
	return nil
}

// AtomicOps returns the atomic operations for this transaction.
func (tx *UnsignedCollectRewardsTx) AtomicOps() (ids.ID, *atomic.Requests, error) {
	// TOOD: I'm leaving the ExportTx implementation here for inspiration
	txID := tx.ID()

	elems := make([]*atomic.Element, 1)
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

	elems[0] = elem

	return tx.DestinationChain, &atomic.Requests{PutRequests: elems}, nil
}

func (vm *VM) NewCollectRewardsTx(
	calculation RewardCalculation,
	blockId ids.ID,
	blockTimestamp uint64,
	coinbase common.Address,
	validatorRewardAddress common.Address,
	incentivePoolRewardAddress common.Address,

) (*Tx, error) {
	// Create the transaction
	utx := &UnsignedCollectRewardsTx{
		NetworkID:                  vm.ctx.NetworkID,
		BlockchainID:               vm.ctx.ChainID,
		DestinationChain:           constants.PlatformChainID,
		BlockId:                    blockId,
		BlockTimestamp:             blockTimestamp,
		RewardCalculation:          calculation,
		Coinbase:                   coinbase,
		ValidatorRewardAddress:     validatorRewardAddress,
		IncentivePoolRewardAddress: incentivePoolRewardAddress,
	}

	// TODO: Make validatorRewardAddress ShortID typed
	pChainAddress, _ := ids.ToShortID(validatorRewardAddress.Bytes())
	utx.ExportedOutputs = []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: vm.ctx.AVAXAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: calculation.ValidatorRewardToExport,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  0,
				Threshold: 1,
				Addrs:     []ids.ShortID{pChainAddress},
			},
		},
	}}

	tx := &Tx{
		UnsignedAtomicTx: utx,
		Creds:            make([]verify.Verifiable, 0),
	}

	log.Info("New CollectionRewardTx\nTo: %s\nAmount: %d\nCurr Coinbase balance: %d\n",
		pChainAddress.String(),
		calculation.ValidatorRewardToExport,
		calculation.PrevValidatorRewards,
	)

	return tx, utx.Verify(vm.ctx, vm.currentRules())
}

// EVMStateTransfer executes the state update from the atomic export transaction
func (tx *UnsignedCollectRewardsTx) EVMStateTransfer(ctx *snow.Context, state *state.StateDB) error {
	state.SubBalance(tx.Coinbase, tx.RewardCalculation.CoinbaseAmountToSub)
	validatorRewards := new(big.Int).Add(tx.RewardCalculation.PrevValidatorRewards, tx.RewardCalculation.ValidatorRewardAmount)
	state.SetState(tx.Coinbase, Slot0, common.BigToHash(validatorRewards))
	state.SetState(tx.Coinbase, Slot1, common.BigToHash(new(big.Int).SetUint64(tx.BlockTimestamp)))

	ipRewards := new(big.Int).Add(tx.RewardCalculation.PrevIncentivePoolRewards, tx.RewardCalculation.IncentivePoolRewardAmount)
	state.AddBalance(tx.IncentivePoolRewardAddress, tx.RewardCalculation.IncentivePoolRewardAmount)
	state.SetState(tx.IncentivePoolRewardAddress, Slot0, common.BigToHash(ipRewards))

	return nil
}
