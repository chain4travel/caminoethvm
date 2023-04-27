package vm

import (
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/coreth/params"
	"github.com/decred/dcrd/dcrec/secp256k1/v3"
)

type caminoRecover struct{}

func (caminoRecover) RequiredGas(input []byte) uint64 {
	dataCopyCost := 2*uint64(len(input)+31)/32*params.IdentityPerWordGas + params.IdentityBaseGas
	hash256Cost := uint64(len(input)+31)/32*params.Sha256PerWordGas + params.Sha256BaseGas
	hash160Cost := uint64(len(input)+31)/32*params.Ripemd160PerWordGas + params.Ripemd160BaseGas
	return dataCopyCost + hash256Cost + hash160Cost
}

// Run implements PrecompiledContract interface. It accepts a 65-byte public key recovered already from the signature.
// It returns the P/X-chain address bytes corresponding to the public key.
func (caminoRecover) Run(input []byte) ([]byte, error) {
	publicKey, err := secp256k1.ParsePubKey(input)
	if err != nil {
		return nil, err
	}
	return hashing.PubkeyBytesToAddress(publicKey.SerializeCompressed()), nil
}
