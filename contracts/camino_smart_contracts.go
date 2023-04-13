package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NextSlot(addr common.Hash) common.Hash {
	bigAddr := new(big.Int).Add(addr.Big(), common.Big1)
	return common.BigToHash(bigAddr)
}

func EntryAddress(address common.Address, slot int64) common.Hash {
	return crypto.Keccak256Hash(address.Hash().Bytes(), common.BigToHash(big.NewInt(slot)).Bytes())
}
