package commands

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/common"
)

var (
	AdminContractAddr = common.HexToAddress("0x010000000000000000000000000000000000000a")

	SharedMemoryCommandBaseID    = ids.ID{0xC, 0xA, 0x1, 0x1, 0xE, 0xD} // CA11ED
	COMMAND_SET_BASE_FEE_OFFSET  = uint64(1)
	COMMAND_SET_KYC_STATE_OFFSET = uint64(2)

	SharedMemoryCommandOffsets = [2]uint64{COMMAND_SET_BASE_FEE_OFFSET, COMMAND_SET_KYC_STATE_OFFSET}

	KYC_OFFSET = uint(1)
	KYB_OFFSET = uint(8)
)
