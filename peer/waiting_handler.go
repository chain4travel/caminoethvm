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

// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"github.com/chain4travel/caminogo/ids"
	"github.com/chain4travel/caminoethvm/plugin/evm/message"
)

var _ message.ResponseHandler = &waitingResponseHandler{}

// waitingResponseHandler implements the ResponseHandler interface
// Internally used to wait for response after making a request synchronously
// responseChan may contain response bytes if the original request has not failed
// responseChan is closed in either fail or success scenario
type waitingResponseHandler struct {
	responseChan chan []byte // blocking channel with response bytes
	failed       bool        // whether the original request is failed
}

// OnResponse passes the response bytes to the responseChan and closes the channel
func (w *waitingResponseHandler) OnResponse(_ ids.ShortID, _ uint32, response []byte) error {
	w.responseChan <- response
	close(w.responseChan)
	return nil
}

// OnFailure sets the failed flag to true and closes the channel
func (w *waitingResponseHandler) OnFailure(ids.ShortID, uint32) error {
	w.failed = true
	close(w.responseChan)
	return nil
}

// newWaitingResponseHandler returns new instance of the waitingResponseHandler
func newWaitingResponseHandler() *waitingResponseHandler {
	return &waitingResponseHandler{responseChan: make(chan []byte)}
}
