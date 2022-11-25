package admin

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	AdminContractAddr = common.HexToAddress("0x010000000000000000000000000000000000000a")
)

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	// SetState(common.Address, common.Hash, common.Hash)

	// SetCode(common.Address, []byte)

	// SetNonce(common.Address, uint64)
	// GetNonce(common.Address) uint64

	// GetBalance(common.Address) *big.Int
	// AddBalance(common.Address, *big.Int)
	// SubBalance(common.Address, *big.Int)

	// CreateAccount(common.Address)
	// Exist(common.Address) bool

	// Suicide(common.Address) bool
	// Finalise(deleteEmptyObjects bool)
}

type KYCStatus int

const (
	KYC_UNKNOWN  = -1
	KYC_VERIFIED = 1
	KYC_EXPIRED  = 2

	FucntionSelectorLength = 4
)

func KycRoleKeyFromAddr(addr common.Address) common.Hash {
	return crypto.Keccak256Hash(append(addr.Hash().Bytes(), common.HexToHash("0x2").Bytes()... /*slot 2 reference admin.sol map(address => uint)*/))
}

// Returns the KYC status for a given address
func GetKYCStatusForAddress(state StateDB, addr common.Address) KYCStatus {
	addrKey := KycRoleKeyFromAddr(addr)
	status := new(big.Int).SetBytes(state.GetState(AdminContractAddr, addrKey).Bytes())
	return KYCStatus(status.Int64())
}

// returns true if a fuction of a SC is blacklist and its execution should be prevented
func IsFunctionBlacklisted(state StateDB, scAddr common.Address, input []byte) error {
	if len(input) >= FucntionSelectorLength {
		var local = make([]byte, len(input)) // directly using the input resulted in execution reverted errors
		copy(local, input)
		funcSignature := local[:FucntionSelectorLength]
		key := crypto.Keccak256Hash(append(funcSignature, scAddr.Bytes()...))
		role := new(big.Int).SetBytes(state.GetState(AdminContractAddr, key).Bytes())
		if role.Cmp(common.Big1) == 0 {
			return fmt.Errorf("function %X on %s is blacklisted", funcSignature, scAddr)
		}
	}
	return nil
}
