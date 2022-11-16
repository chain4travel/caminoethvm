package precompile_interface

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	AdminContractAddr = common.HexToAddress("0x010000000000000000000000000000000000000a")
	KYC_APPROVED      = common.HexToHash("0x01")
	KYC_EXPIRED       = common.HexToHash("0x02")
)

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	SetCode(common.Address, []byte)

	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64

	GetBalance(common.Address) *big.Int
	AddBalance(common.Address, *big.Int)
	SubBalance(common.Address, *big.Int)

	CreateAccount(common.Address)
	Exist(common.Address) bool

	Suicide(common.Address) bool
	// Finalise(deleteEmptyObjects bool)
}

type KYC_STATUS int

const (
	KYC_STATUS_UNKNOWN  = -1
	KYC_STATUS_APPROVED = 1
	KYC_STATUS_EXPIRED  = 2
)

func GetKYCStatusForAddress(stateDB StateDB, addr common.Address) KYC_STATUS {
	addrKey := common.BytesToHash(crypto.Keccak256(append(addr.Hash().Bytes(), common.HexToHash("0x2").Bytes()... /*slot 2 referenece admin.sol map(address => uint)*/)))
	roleRaw := stateDB.GetState(AdminContractAddr, addrKey)
	switch roleRaw.String() {
	case KYC_APPROVED.String():
		return KYC_STATUS_APPROVED
	case KYC_EXPIRED.String():
		return KYC_STATUS_EXPIRED
	default:
		return KYC_STATUS_UNKNOWN
	}
}
