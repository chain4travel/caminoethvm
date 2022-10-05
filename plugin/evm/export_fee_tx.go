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
	"math/big"

	"github.com/chain4travel/caminoethvm/core/state"
	"github.com/chain4travel/caminoethvm/params"

	"github.com/chain4travel/caminogo/chains/atomic"
	"github.com/chain4travel/caminogo/ids"
	"github.com/chain4travel/caminogo/snow"
	"github.com/chain4travel/caminogo/vms/components/avax"
)

// UnsignedExportFeeTx is an internal fee collection transaction
type UnsignedExportFeeTx struct {
	avax.Metadata
	// ID of the network on which this tx was issued
	NetworkID uint32 `serialize:"true" json:"networkID"`
	// ID of this blockchain.
	BlockchainID ids.ID `serialize:"true" json:"blockchainID"`
	// Which chain to send the funds to
	DestinationChain ids.ID `serialize:"true" json:"destinationChain"`
	// Outputs that are exported to the chain
	ExportedOutputs []*avax.TransferableOutput `serialize:"true" json:"exportedOutputs"`

	// TODO: Fee collection specifics:
	// - Results of the fee calculation (reward & incentive pool)
	// - Block timestamp at which the fee was calculated
	// My idea is to prepare such Tx at block's Accept() with all needed data
	// and then add the Tx to mempool.
}

// InputUTXOs returns a set of all the hash(address:nonce) exporting funds.
func (tx *UnsignedExportFeeTx) InputUTXOs() ids.Set {
	// Not sure it will be needed - mock
	return ids.NewSet(0)
}

// Verify this transaction is well-formed
func (tx *UnsignedExportFeeTx) Verify(
	ctx *snow.Context,
	rules params.Rules,
) error {
	return nil
}

func (tx *UnsignedExportFeeTx) GasUsed(bool) (uint64, error) {
	return 0, nil
}

// Amount of [assetID] burned by this transaction
func (tx *UnsignedExportFeeTx) Burned(assetID ids.ID) (uint64, error) {
	// Not sure it will be needed - mock
	return 0, nil
}

// SemanticVerify this transaction is valid.
func (tx *UnsignedExportFeeTx) SemanticVerify(
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
func (tx *UnsignedExportFeeTx) AtomicOps() (ids.ID, *atomic.Requests, error) {
	// TOOD: I'm leaving the ExportTx implementation here for inspiration
	txID := tx.ID()

	elems := make([]*atomic.Element, len(tx.ExportedOutputs))
	for i, out := range tx.ExportedOutputs {
		utxo := &avax.UTXO{
			UTXOID: avax.UTXOID{
				TxID:        txID,
				OutputIndex: uint32(i),
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

		elems[i] = elem
	}

	return tx.DestinationChain, &atomic.Requests{PutRequests: elems}, nil
}

func (vm *VM) newExportFeeTx(
	assetID ids.ID, // AssetID of the tokens to export
	amount uint64, // Amount of tokens to export
	chainID ids.ID, // Chain to send the UTXOs to
	to ids.ShortID, // Address of chain recipient
	// TODO: other params
) (*Tx, error) {
	// Create the transaction
	utx := &UnsignedExportFeeTx{
		NetworkID:        vm.ctx.NetworkID,
		BlockchainID:     vm.ctx.ChainID,
		DestinationChain: chainID,
		ExportedOutputs:  []*avax.TransferableOutput{},
	}
	tx := &Tx{
		UnsignedAtomicTx: utx,
		Creds:            make([]verify.Verifiable, 0),
	}

	return tx, utx.Verify(vm.ctx, vm.currentRules())
}

// EVMStateTransfer executes the state update from the atomic export transaction
func (tx *UnsignedExportFeeTx) EVMStateTransfer(ctx *snow.Context, state *state.StateDB) error {
	// here is a place to update the state
	// e.g. state.setState(BlackHoleaddress, slot0, amount)
	return nil
}
