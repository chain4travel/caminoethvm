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
	"context"
	"encoding/json"

	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/chain4travel/caminoethvm/core"
)

// StaticService defines the static API services exposed by the evm
type StaticService struct{}

// BuildGenesisReply is the reply from BuildGenesis
type BuildGenesisReply struct {
	Bytes    string              `json:"bytes"`
	Encoding formatting.Encoding `json:"encoding"`
}

// BuildGenesis returns the UTXOs such that at least one address in [args.Addresses] is
// referenced in the UTXO.
func (*StaticService) BuildGenesis(_ context.Context, args *core.Genesis) (*BuildGenesisReply, error) {
	bytes, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	bytesStr, err := formatting.EncodeWithChecksum(formatting.Hex, bytes)
	if err != nil {
		return nil, err
	}
	return &BuildGenesisReply{
		Bytes:    bytesStr,
		Encoding: formatting.Hex,
	}, nil
}
