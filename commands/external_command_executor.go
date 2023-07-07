package commands

import (
	"math/big"

	"github.com/ava-labs/coreth/core/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	_ ExternalCommandVisitor = (*ExternalCommandExecutor)(nil)
)

type ExternalCommandExecutor struct {
	stateDB *state.StateDB
}

func NewExternalCommandExecutor(stateDB *state.StateDB) (ExternalCommandExecutor, error) {
	return ExternalCommandExecutor{
		stateDB,
	}, nil
}

func (ece ExternalCommandExecutor) ExecuteSetBaseFeeCommand(cmd *ExternalCommandSetBaseFee) error {
	ece.stateDB.SetState(AdminContractAddr, crypto.Keccak256Hash(common.HexToHash("0x1").Bytes()), common.BigToHash(&cmd.NewBaseFee))
	return nil
}

func (ece ExternalCommandExecutor) ExecuteSetKYCStateCommand(cmd *ExternalCommandSetKYCState) error {
	for _, update := range cmd.KYCUpdates {
		kycState := big.NewInt(0)
		if update.KYCVerified {
			kycState = kycState.Or(kycState, common.Big1.Lsh(common.Big1, KYC_OFFSET))
		}

		if update.KYBVerified {
			kycState = kycState.Or(kycState, common.Big1.Lsh(common.Big1, KYB_OFFSET))
		}

		ece.stateDB.SetState(AdminContractAddr, KycRoleKeyFromAddr(update.Address), common.BigToHash(kycState))
	}

	return nil
}

func KycRoleKeyFromAddr(addr common.Address) common.Hash {
	return crypto.Keccak256Hash(addr.Hash().Bytes(), common.HexToHash("0x2").Bytes() /*slot 2 reference admin.sol map(address => uint)*/)
}
