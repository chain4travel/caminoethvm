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
	"context"
	"time"

	"github.com/chain4travel/caminogo/codec"
	"github.com/chain4travel/caminogo/ids"

	"github.com/chain4travel/caminoethvm/core/rawdb"
	"github.com/chain4travel/caminoethvm/ethdb"
	"github.com/chain4travel/caminoethvm/plugin/evm/message"
	"github.com/chain4travel/caminoethvm/statesync/handlers/stats"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// CodeRequestHandler is a peer.RequestHandler for message.CodeRequest
// serving requested contract code bytes
type CodeRequestHandler struct {
	codeReader ethdb.KeyValueReader
	codec      codec.Manager
	stats      stats.HandlerStats
}

func NewCodeRequestHandler(codeReader ethdb.KeyValueReader, stats stats.HandlerStats, codec codec.Manager) *CodeRequestHandler {
	handler := &CodeRequestHandler{
		codeReader: codeReader,
		codec:      codec,
		stats:      stats,
	}
	return handler
}

// OnCodeRequest handles request to retrieve contract code by its hash in message.CodeRequest
// Never returns error
// Returns nothing if code hash is not found
// Expects returned errors to be treated as FATAL
// Assumes ctx is active
func (n *CodeRequestHandler) OnCodeRequest(_ context.Context, nodeID ids.ShortID, requestID uint32, codeRequest message.CodeRequest) ([]byte, error) {
	startTime := time.Now()
	n.stats.IncCodeRequest()

	// always report code read time metric
	defer func() {
		n.stats.UpdateCodeReadTime(time.Since(startTime))
	}()

	codeData := rawdb.ReadCode(n.codeReader, codeRequest.Hash)
	if len(codeData) == 0 {
		n.stats.IncMissingCodeHash()
		log.Debug("requested code not found, dropping request", "nodeID", nodeID, "requestID", requestID, "hash", codeRequest.Hash)
		return nil, nil
	}

	codeResponse := message.CodeResponse{Data: common.CopyBytes(codeData)}
	responseBytes, err := n.codec.Marshal(message.Version, codeResponse)
	if err != nil {
		log.Warn("could not marshal CodeResponse, dropping request", "nodeID", nodeID, "requestID", requestID, "hash", codeRequest.Hash, "err", err)
		return nil, nil
	}
	n.stats.UpdateCodeBytesReturned(uint32(len(codeData)))
	return responseBytes, nil
}
