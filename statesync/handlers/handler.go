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

	"github.com/chain4travel/caminogo/ids"

	"github.com/chain4travel/caminoethvm/plugin/evm/message"
)

var _ message.RequestHandler = &syncHandler{}

type syncHandler struct {
	stateTrieLeafsRequestHandler  *LeafsRequestHandler
	atomicTrieLeafsRequestHandler *LeafsRequestHandler
	blockRequestHandler           *BlockRequestHandler
	codeRequestHandler            *CodeRequestHandler
}

func NewSyncHandler(
	stateTrieLeafsRequestHandler *LeafsRequestHandler,
	atomicTrieLeafsRequestHandler *LeafsRequestHandler,
	blockRequestHandler *BlockRequestHandler,
	codeRequestHandler *CodeRequestHandler,
) message.RequestHandler {
	return &syncHandler{
		stateTrieLeafsRequestHandler:  stateTrieLeafsRequestHandler,
		atomicTrieLeafsRequestHandler: atomicTrieLeafsRequestHandler,
		blockRequestHandler:           blockRequestHandler,
		codeRequestHandler:            codeRequestHandler,
	}
}

func (s *syncHandler) HandleStateTrieLeafsRequest(ctx context.Context, nodeID ids.ShortID, requestID uint32, leafsRequest message.LeafsRequest) ([]byte, error) {
	return s.stateTrieLeafsRequestHandler.OnLeafsRequest(ctx, nodeID, requestID, leafsRequest)
}

func (s *syncHandler) HandleAtomicTrieLeafsRequest(ctx context.Context, nodeID ids.ShortID, requestID uint32, leafsRequest message.LeafsRequest) ([]byte, error) {
	return s.atomicTrieLeafsRequestHandler.OnLeafsRequest(ctx, nodeID, requestID, leafsRequest)
}

func (s *syncHandler) HandleBlockRequest(ctx context.Context, nodeID ids.ShortID, requestID uint32, blockRequest message.BlockRequest) ([]byte, error) {
	return s.blockRequestHandler.OnBlockRequest(ctx, nodeID, requestID, blockRequest)
}

func (s *syncHandler) HandleCodeRequest(ctx context.Context, nodeID ids.ShortID, requestID uint32, codeRequest message.CodeRequest) ([]byte, error) {
	return s.codeRequestHandler.OnCodeRequest(ctx, nodeID, requestID, codeRequest)
}
