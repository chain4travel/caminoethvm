package commands

import (
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/components/verify"
	"github.com/ethereum/go-ethereum/common"
)

type ExternalCommandVisitor interface {
	ExecuteSetBaseFeeCommand(*ExternalCommandSetBaseFee) error
	ExecuteSetKYCStateCommand(*ExternalCommandSetKYCState) error
}

type ExternalCommand interface {
	verify.Verifiable

	Visit(ExternalCommandVisitor) error
	SharedMemoryID() ids.ID
}

type ExternalCommandSetBaseFee struct {
	ExternalCommand `serializable:"true"`
	NewBaseFee      big.Int `serializable:"true"`
}

type ExternalCommandSetKYCState struct {
	ExternalCommand `serializable:"true"`
	KYCUpdates      []KYCUpdate `serialize:"true", json:"kyc_updates"`
}

type KYCUpdate struct {
	Address     common.Address `serialize:"true", json:"address"`
	KYCVerified bool           `serialize:"true", json:"kyc_verified"`
	KYBVerified bool           `serialize:"true", json:"kyb_verified"`
}

func (cmd *ExternalCommandSetBaseFee) Visit(visitor ExternalCommandVisitor) error {
	return visitor.ExecuteSetBaseFeeCommand(cmd)
}

func (cmd *ExternalCommandSetBaseFee) SharedMemoryOffset() ids.ID {
	return SharedMemoryCommandBaseID.Prefix(COMMAND_SET_BASE_FEE_OFFSET)
}

func (cmd *ExternalCommandSetKYCState) Visit(visitor ExternalCommandVisitor) error {
	return visitor.ExecuteSetKYCStateCommand(cmd)
}

func (cmd *ExternalCommandSetKYCState) SharedMemoryOffset() ids.ID {
	return SharedMemoryCommandBaseID.Prefix(COMMAND_SET_KYC_STATE_OFFSET)
}
