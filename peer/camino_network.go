// Copyright (C) 2022-2023, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/sync/semaphore"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/engine/common"

	"github.com/ava-labs/coreth/peer/stats"
	"github.com/ava-labs/coreth/plugin/evm/message"
)

type caminoNetwork struct {
	network
}

func NewCaminoNetwork(appSender common.AppSender, codec codec.Manager, crossChainCodec codec.Manager, self ids.NodeID, maxActiveAppRequests int64, maxActiveCrossChainRequests int64) Network {
	return &caminoNetwork{
		network: network{
			appSender:                  appSender,
			codec:                      codec,
			crossChainCodec:            crossChainCodec,
			self:                       self,
			outstandingRequestHandlers: make(map[uint32]message.ResponseHandler),
			activeAppRequests:          semaphore.NewWeighted(maxActiveAppRequests),
			activeCrossChainRequests:   semaphore.NewWeighted(maxActiveCrossChainRequests),
			gossipHandler:              message.NoopMempoolGossipHandler{},
			appRequestHandler:          message.NoopRequestHandler{},
			crossChainRequestHandler:   message.NoopCrossChainRequestHandler{},
			peers:                      NewPeerTracker(),
			appStats:                   stats.NewRequestHandlerStats(),
			crossChainStats:            stats.NewCrossChainRequestHandlerStats(),
		},
	}
}

// CrossChainAppRequest notifies the VM when another chain in the network requests for data.
// Send a CrossChainAppResponse to [chainID] in response to a valid message using the same
// [requestID] before the deadline.
func (cn *caminoNetwork) CrossChainAppRequest(ctx context.Context, requestingChainID ids.ID, requestID uint32, deadline time.Time, request []byte) error {
	return cn.network.CrossChainAppRequest(ctx, requestingChainID, requestID, deadline, request)
}

func (n *network) RequestCrossChain(chainID ids.ID, msg []byte, handler message.ResponseHandler) error {
	if handler != nil {
		return fmt.Errorf("ResponseHandler not yet supported")
	}
	log.Debug("sending request to chain", "chainID", chainID, "requestLen", len(msg))
	return n.appSender.SendCrossChainAppRequest(context.TODO(), chainID, 0, msg)
}
