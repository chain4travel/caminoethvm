// Copyright (C) 2022, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"github.com/chain4travel/caminoethvm/core/state"
)

// Gas Price
const (
	SunrisePhase0ExtraDataSize = 0

//	SunrisePhase0BaseFee     uint64 = 50_000_000_000
)

func SunrisePhase0BaseFee() uint64 {
	var s *state.StateDB

	basefee := s.GetBaseFee().Uint64()
	return basefee

}
