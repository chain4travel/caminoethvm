// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/params"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/txs/commands"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
)

var (
	_                            UnsignedAtomicTx       = &UnsignedExecuteExternalCommandTx{}
	_                            secp256k1fx.UnsignedTx = &UnsignedExecuteExternalCommandTx{}
	errOnlyPChainCommandsAllowed                        = errors.New("only commands sent from P-Chain are allowed to be processed")
	errCommandParseError                                = errors.New("failed to parse command from bytes")
	errFailedToConstructTx                              = errors.New("failed to construct tx")
)

// UnsignedExecuteExternalCommandTx
// Loads the a Command from SharedMemory, validates and executes it against the EVM Database
type UnsignedExecuteExternalCommandTx struct {
	avax.Metadata
	// ID of the network on which this tx was issued
	NetworkID uint32 `serialize:"true" json:"networkID"`
	// ID of this blockchain.
	BlockchainID ids.ID `serialize:"true" json:"blockchainID"`
	// Futureproof to recieve potential commands from other chains that P-Chain
	SourceChainID ids.ID `serialize:"true" json:"sourceChainID"`
	// Which command in the SharedMemory to comsume
	//? No particular reason those could not be an array.
	ExternalCommandTxID ids.ID `serialize:"true" json:"externalCommandTxID"`
	// Not sure if this is nessecary but it decouples execution from shared mem
	ExternalCommandBytes []byte `serialize:"true" json:"externalCommandBytes"`

	externalCommand commands.ExternalCommand
}

// InputUTXOs returns the UTXOIDs of the imported funds
func (utx *UnsignedExecuteExternalCommandTx) InputUTXOs() set.Set[ids.ID] {
	return set.NewSet[ids.ID](0)
}

// Verify this transaction is well-formed
func (utx *UnsignedExecuteExternalCommandTx) Verify(
	ctx *snow.Context,
	rules params.Rules,
) error {
	switch {
	case utx == nil:
		return errNilTx
	case utx.externalCommand != nil:
		return fmt.Errorf("pre popluating external Command field is not allowed")
	case utx.NetworkID != ctx.NetworkID:
		return errWrongNetworkID
	case ctx.ChainID != utx.BlockchainID:
		return errWrongBlockchainID
	case utx.SourceChainID != constants.PlatformChainID:
		return errOnlyPChainCommandsAllowed

	}

	return nil
}

func (utx *UnsignedExecuteExternalCommandTx) GasUsed(fixedFee bool) (uint64, error) {

	return 0, nil
}

// Amount of [assetID] burned by this transaction
func (utx *UnsignedExecuteExternalCommandTx) Burned(assetID ids.ID) (uint64, error) {
	return 0, nil
}

// SemanticVerify this transaction is valid.
func (utx *UnsignedExecuteExternalCommandTx) SemanticVerify(
	vm *VM,
	stx *Tx,
	parent *Block,
	baseFee *big.Int,
	rules params.Rules,
) error {
	if err := utx.Verify(vm.ctx, rules); err != nil {
		return err
	}

	if !vm.bootstrapped {
		// Allow for force committing during bootstrapping
		return nil
	}

	// check for existence inside of shared Memory
	allCommandBytes, err := vm.ctx.SharedMemory.Get(utx.SourceChainID, [][]byte{utx.ExternalCommandTxID[:]})
	if err != nil {
		return fmt.Errorf("failed to fetch import UTXOs from %s due to: %w", utx.SourceChainID, err)
	}

	if len(allCommandBytes) != 0 {
		return fmt.Errorf("zero length command in shared memory @ %s", utx.ExternalCommandTxID)
	}

	if !bytes.Equal(utx.ExternalCommandBytes, allCommandBytes[0]) {
		return fmt.Errorf("missmatched command in shared memory @ %s", utx.ExternalCommandTxID)
	}

	// I know this sucks, but I havent had a better idea, i have no vm inside of the EVM StateTransfer
	_, err = commands.Codec.Unmarshal(allCommandBytes[0], &utx.externalCommand)
	if err != nil {
		return fmt.Errorf("failed to unmarshal external command: %v", err)
	}

	return vm.conflicts(utx.InputUTXOs(), parent)
}

// AtomicOps returns imported inputs spent on this transaction
// We spend imported UTXOs here rather than in semanticVerify because
// we don't want to remove an imported UTXO in semanticVerify
// only to have the transaction not be Accepted. This would be inconsistent.
// Recall that imported UTXOs are not kept in a versionDB.
func (utx *UnsignedExecuteExternalCommandTx) AtomicOps() (ids.ID, *atomic.Requests, error) {
	return utx.SourceChainID, &atomic.Requests{RemoveRequests: [][]byte{utx.ExternalCommandTxID[:]}}, nil
}

// newExecuteExternalCommandTx returns a new ExecuteExternalCommandTx
func (vm *VM) newExecuteExternalCommandTx(
	externalCommandTxID ids.ID,
) (*Tx, error) {

	externalCommandBytes, err := vm.ctx.SharedMemory.Get(constants.PlatformChainID, [][]byte{externalCommandTxID[:]})
	if err != nil {
		return nil, fmt.Errorf("problem retrieving atomic UTXOs: %w", err)
	}

	utx := &UnsignedExecuteExternalCommandTx{
		NetworkID:            vm.ctx.NetworkID,
		BlockchainID:         vm.ctx.ChainID,
		SourceChainID:        constants.PlatformChainID,
		ExternalCommandTxID:  externalCommandTxID,
		ExternalCommandBytes: externalCommandBytes[0],
	}
	tx := &Tx{UnsignedAtomicTx: utx}
	if err := tx.Sign(vm.codec, nil); err != nil {
		return nil, err
	}

	return tx, utx.Verify(vm.ctx, vm.currentRules())
}

// EVMStateTransfer performs the state transfer to increase the balances of
// accounts accordingly with the imported EVMOutputs
// Invarriant: SemanticVerify has to be executed before this
func (utx *UnsignedExecuteExternalCommandTx) EVMStateTransfer(ctx *snow.Context, state *state.StateDB) error {
	if utx.externalCommand == nil {
		return fmt.Errorf("FATAL: Invariant violated - SemanticVerify not executed")
	}

	return utx.externalCommand.EVMStateTransfer(ctx, state)
}

func (vm *VM) processCrossChainCommandMessage(msg []byte) error {
	var cmdMsg message.CaminoCommandMessage
	_, err := message.Codec.Unmarshal(msg, &cmdMsg)
	if err != nil {
		return err
	}

	tx, err := vm.newExecuteExternalCommandTx(cmdMsg.CommandTxID)
	if err != nil {
		return err
	}

	return vm.issueTx(tx, true)
}
