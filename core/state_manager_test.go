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

package core

import (
	"math/big"
	"testing"

	"github.com/chain4travel/caminoethvm/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type MockTrieDB struct {
	LastReference   common.Hash
	LastDereference common.Hash
	LastCommit      common.Hash
}

func (t *MockTrieDB) Reference(child common.Hash, parent common.Hash) {
	t.LastReference = child
}
func (t *MockTrieDB) Dereference(root common.Hash) {
	t.LastDereference = root
}
func (t *MockTrieDB) Commit(root common.Hash, report bool, callback func(common.Hash)) error {
	t.LastCommit = root
	return nil
}
func (t *MockTrieDB) Size() (common.StorageSize, common.StorageSize) {
	return 0, 0
}
func (t *MockTrieDB) Cap(limit common.StorageSize) error {
	return nil
}

func TestCappedMemoryTrieWriter(t *testing.T) {
	m := &MockTrieDB{}
	w := NewTrieWriter(m, &CacheConfig{Pruning: true})
	assert := assert.New(t)
	for i := 0; i < commitInterval+1; i++ {
		bigI := big.NewInt(int64(i))
		block := types.NewBlock(
			&types.Header{
				Root:   common.BigToHash(bigI),
				Number: bigI,
			},
			nil, nil, nil, nil, nil, true,
		)

		assert.NoError(w.InsertTrie(block))
		assert.Equal(block.Root(), m.LastReference, "should not have referenced block on insert")
		assert.Equal(common.Hash{}, m.LastDereference, "should not have dereferenced block on insert")
		assert.Equal(common.Hash{}, m.LastCommit, "should not have committed block on insert")
		m.LastReference = common.Hash{}

		w.AcceptTrie(block)
		assert.Equal(common.Hash{}, m.LastReference, "should not have referenced block on accept")
		if i < tipBufferSize {
			assert.Equal(common.Hash{}, m.LastDereference, "should not have dereferenced block on accept")
		} else {
			assert.Equal(common.BigToHash(big.NewInt(int64(i-tipBufferSize))), m.LastDereference, "should have dereferenced old block on last accept")
			m.LastDereference = common.Hash{}
		}
		if i < commitInterval {
			assert.Equal(common.Hash{}, m.LastCommit, "should not have committed block on accept")
		} else {
			assert.Equal(block.Root(), m.LastCommit, "should have committed block after commitInterval")
			m.LastCommit = common.Hash{}
		}

		w.RejectTrie(block)
		assert.Equal(common.Hash{}, m.LastReference, "should not have referenced block on reject")
		assert.Equal(block.Root(), m.LastDereference, "should have dereferenced block on reject")
		assert.Equal(common.Hash{}, m.LastCommit, "should not have committed block on reject")
		m.LastDereference = common.Hash{}
	}
}

func TestNoPruningTrieWriter(t *testing.T) {
	m := &MockTrieDB{}
	w := NewTrieWriter(m, &CacheConfig{})
	assert := assert.New(t)
	for i := 0; i < tipBufferSize+1; i++ {
		bigI := big.NewInt(int64(i))
		block := types.NewBlock(
			&types.Header{
				Root:   common.BigToHash(bigI),
				Number: bigI,
			},
			nil, nil, nil, nil, nil, true,
		)

		assert.NoError(w.InsertTrie(block))
		assert.Equal(common.Hash{}, m.LastReference, "should not have referenced block on insert")
		assert.Equal(common.Hash{}, m.LastDereference, "should not have dereferenced block on insert")
		assert.Equal(block.Root(), m.LastCommit, "should have committed block on insert")
		m.LastCommit = common.Hash{}

		w.AcceptTrie(block)
		assert.Equal(common.Hash{}, m.LastReference, "should not have referenced block on accept")
		assert.Equal(common.Hash{}, m.LastDereference, "should not have dereferenced block on accept")
		assert.Equal(common.Hash{}, m.LastCommit, "should not have committed block on accept")

		w.RejectTrie(block)
		assert.Equal(common.Hash{}, m.LastReference, "should not have referenced block on reject")
		assert.Equal(common.Hash{}, m.LastDereference, "should not have dereferenced block on reject")
		assert.Equal(common.Hash{}, m.LastCommit, "should not have committed block on reject")
	}
}
