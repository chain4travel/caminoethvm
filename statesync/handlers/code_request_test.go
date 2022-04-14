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

// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"

	"github.com/chain4travel/caminoethvm/params"

	"github.com/chain4travel/caminogo/ids"
	"github.com/chain4travel/caminoethvm/core/rawdb"
	"github.com/chain4travel/caminoethvm/ethdb/memorydb"
	"github.com/chain4travel/caminoethvm/plugin/evm/message"
	"github.com/chain4travel/caminoethvm/statesync/handlers/stats"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCodeRequestHandler(t *testing.T) {
	codec, err := message.BuildCodec()
	if err != nil {
		t.Fatal("unexpected error when building codec", err)
	}

	database := memorydb.New()

	codeBytes := []byte("some code goes here")
	codeHash := crypto.Keccak256Hash(codeBytes)
	rawdb.WriteCode(database, codeHash, codeBytes)

	codeRequestHandler := NewCodeRequestHandler(database, stats.NewNoopHandlerStats(), codec)

	// query for known code entry
	responseBytes, err := codeRequestHandler.OnCodeRequest(context.Background(), ids.GenerateTestShortID(), 1, message.CodeRequest{Hash: codeHash})
	assert.NoError(t, err)

	var response message.CodeResponse
	if _, err = codec.Unmarshal(responseBytes, &response); err != nil {
		t.Fatal("error unmarshalling CodeResponse", err)
	}
	assert.True(t, bytes.Equal(codeBytes, response.Data))

	// query for missing code entry
	responseBytes, err = codeRequestHandler.OnCodeRequest(context.Background(), ids.GenerateTestShortID(), 2, message.CodeRequest{Hash: common.BytesToHash([]byte("some unknown hash"))})
	assert.NoError(t, err)
	assert.Nil(t, responseBytes)

	// assert max size code bytes are handled
	codeBytes = make([]byte, params.MaxCodeSize)
	n, err := rand.Read(codeBytes)
	assert.NoError(t, err)
	assert.Equal(t, params.MaxCodeSize, n)
	codeHash = crypto.Keccak256Hash(codeBytes)
	rawdb.WriteCode(database, codeHash, codeBytes)

	responseBytes, err = codeRequestHandler.OnCodeRequest(context.Background(), ids.GenerateTestShortID(), 3, message.CodeRequest{Hash: codeHash})
	assert.NoError(t, err)
	assert.NotNil(t, responseBytes)

	response = message.CodeResponse{}
	if _, err = codec.Unmarshal(responseBytes, &response); err != nil {
		t.Fatal("error unmarshalling CodeResponse", err)
	}
	assert.True(t, bytes.Equal(codeBytes, response.Data))
}
