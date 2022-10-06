// Copyright (C) 2022, Chain4Travel AG. All rights reserved.

package admin

import (
	"math/big"

	"github.com/chain4travel/caminoethvm/core/state"
	"github.com/chain4travel/caminoethvm/core/types"
)

// Admin interface to control administrative tasks, which are intended to
// be controlled on block level like BaseFee or Blacklisting or KYC
type AdminController interface {
	// Get the FixedBaseFee which should applied for blocks after height
	GetFixedBaseFee(head *types.Header, state *state.StateDB) *big.Int
}
