// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package dummy

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/coreth/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type OnFinalizeCallbackType = func(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header)
type OnFinalizeAndAssembleCallbackType = func(header *types.Header, state *state.StateDB, txs []*types.Transaction) ([]byte, error)
type OnAPIsCallbackType = func(consensus.ChainHeaderReader) []rpc.API
type OnExtraStateChangeType = func(block *types.Block, statedb *state.StateDB) error

type ConsensusCallbacks struct {
	OnAPIs                OnAPIsCallbackType
	OnFinalize            OnFinalizeCallbackType
	OnFinalizeAndAssemble OnFinalizeAndAssembleCallbackType
	OnExtraStateChange    OnExtraStateChangeType
}

type DummyEngine struct {
	cb *ConsensusCallbacks
}

func NewDummyEngine(cb *ConsensusCallbacks) *DummyEngine {
	return &DummyEngine{
		cb: cb,
	}
}

func NewFaker() *DummyEngine {
	return NewDummyEngine(new(ConsensusCallbacks))
}

var (
	allowedFutureBlockTime = 10 * time.Second // Max time from current time allowed for blocks, before they're considered future blocks
)

var (
	errInvalidBlockTime  = errors.New("timestamp less than parent's")
	errUnclesUnsupported = errors.New("uncles unsupported")
)

// modified from consensus.go
func (self *DummyEngine) verifyHeader(chain consensus.ChainHeaderReader, header, parent *types.Header, uncle bool) error {
	// Ensure that we do not verify an uncle
	if uncle {
		return errUnclesUnsupported
	}
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if header.Time > uint64(time.Now().Add(allowedFutureBlockTime).Unix()) {
		return consensus.ErrFutureBlock
	}
	//if header.Time <= parent.Time {
	if header.Time < parent.Time {
		return errInvalidBlockTime
	}
	// Verify that the gas limit is <= 2^63-1
	cap := uint64(0x7fffffffffffffff)
	if header.GasLimit > cap {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, cap)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	// TODO(aaronbuchwald) verify gas limit EIP-1559
	// TODO(aaronbuchwald) verify base fee is not there pre Apricot Phase 4 and is correct post Apricot Phase 4
	if config := chain.Config(); config.IsApricotPhase1(new(big.Int).SetUint64((header.Time))) {
		if header.GasLimit != params.ApricotPhase1GasLimit {
			return fmt.Errorf("expected gas limit to be %d, but found %d", params.ApricotPhase1GasLimit, header.GasLimit)
		}
	} else {
		// Verify that the gas limit remains within allowed bounds
		diff := int64(parent.GasLimit) - int64(header.GasLimit)
		if diff < 0 {
			diff *= -1
		}
		limit := parent.GasLimit / params.GasLimitBoundDivisor

		if uint64(diff) >= limit || header.GasLimit < params.MinGasLimit {
			return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
		}
	}

	// if !chain.Config().IsLondon(header.Number) {
	// 	// Verify BaseFee not present before EIP-1559 fork.
	// 	if header.BaseFee != nil {
	// 		return fmt.Errorf("invalid baseFee before fork: have %d, want <nil>", header.BaseFee)
	// 	}
	// 	if err := misc.VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
	// 		return err
	// 	}
	// } else if err := misc.VerifyEip1559Header(chain.Config(), parent, header); err != nil {
	// 	// Verify the header's EIP-1559 attributes.
	// 	return err
	// }

	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	// Verify the engine specific seal securing the block
	return self.VerifySeal(chain, header)
}

func (self *DummyEngine) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (self *DummyEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Short circuit if the header is known, or it's parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return self.verifyHeader(chain, header, parent, false)
}

func (self *DummyEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errUnclesUnsupported
	}
	return nil
}

func (self *DummyEngine) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	return nil
}

func (self *DummyEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Difficulty = big.NewInt(1)
	return nil
}

func (self *DummyEngine) Finalize(
	chain consensus.ChainHeaderReader, header *types.Header,
	state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	if self.cb.OnFinalize != nil {
		self.cb.OnFinalize(chain, header, state, txs, uncles)
	}
}

func (self *DummyEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	var extdata []byte
	if self.cb.OnFinalizeAndAssemble != nil {
		ret, err := self.cb.OnFinalizeAndAssemble(header, state, txs)
		extdata = ret
		if err != nil {
			return nil, err
		}
	}
	// commit the final state root
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))

	// Header seems complete, assemble into a block and return
	return types.NewBlock(
		header, txs, uncles, receipts, new(trie.Trie), extdata,
		chain.Config().IsApricotPhase1(new(big.Int).SetUint64(header.Time)),
	), nil
}

func (self *DummyEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

func (self *DummyEngine) APIs(chain consensus.ChainHeaderReader) (res []rpc.API) {
	res = nil
	if self.cb.OnAPIs != nil {
		res = self.cb.OnAPIs(chain)
	}
	return
}

func (self *DummyEngine) Close() error {
	return nil
}

func (self *DummyEngine) ExtraStateChange(block *types.Block, statedb *state.StateDB) error {
	if self.cb.OnExtraStateChange != nil {
		return self.cb.OnExtraStateChange(block, statedb)
	}
	return nil
}

// TODO(aaronbuchwald) ensure that baseFee is correctly serialized on block headers as this is handled in other consensus engines
// to some extent - the seal hash
