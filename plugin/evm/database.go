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
	"github.com/ava-labs/avalanchego/database"
	"github.com/chain4travel/caminoethvm/ethdb"
)

// Database implements ethdb.Database
type Database struct{ database.Database }

// NewBatch implements ethdb.Database
func (db Database) NewBatch() ethdb.Batch { return Batch{db.Database.NewBatch()} }

// NewIterator implements ethdb.Database
//
// Note: This method assumes that the prefix is NOT part of the start, so there's
// no need for the caller to prepend the prefix to the start.
func (db Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	// avalanchego's database implementation assumes that the prefix is part of the
	// start, so it is added here (if it is provided).
	if len(prefix) > 0 {
		newStart := make([]byte, len(prefix)+len(start))
		copy(newStart, prefix)
		copy(newStart[len(prefix):], start)
		start = newStart
	}
	return db.Database.NewIteratorWithStartAndPrefix(start, prefix)
}

// NewIteratorWithStart implements ethdb.Database
func (db Database) NewIteratorWithStart(start []byte) ethdb.Iterator {
	return db.Database.NewIteratorWithStart(start)
}

// Batch implements ethdb.Batch
type Batch struct{ database.Batch }

// ValueSize implements ethdb.Batch
func (batch Batch) ValueSize() int { return batch.Batch.Size() }

// Replay implements ethdb.Batch
func (batch Batch) Replay(w ethdb.KeyValueWriter) error { return batch.Batch.Replay(w) }
