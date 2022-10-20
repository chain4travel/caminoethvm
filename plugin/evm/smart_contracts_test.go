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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/chain4travel/caminoethvm/contracts/build"
	"github.com/chain4travel/caminoethvm/core/rawdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	c4tBind "github.com/chain4travel/caminoethvm/accounts/abi/bind"
	c4tBackends "github.com/chain4travel/caminoethvm/accounts/abi/bind/backends"
	c4tCore "github.com/chain4travel/caminoethvm/core"
)

func TestDummy(t *testing.T) {
	conn, _, auth, _ := newMockEthConnection(uint64(10000000), big.NewInt(9223372036854775807))
	defer func(conn *c4tBackends.SimulatedBackend) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	contractAddr := common.HexToAddress("0x010000000000000000000000000000000000000a")
	newKey, _ := crypto.GenerateKey()
	newPublicKey := newKey.Public()
	publicKeyECDSA, _ := newPublicKey.(*ecdsa.PublicKey)
	newAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	contract, _ := build.NewBuild(contractAddr, conn)
	session := build.BuildSession{Contract: contract, TransactOpts: *auth}
	conn.Commit(true)

	roles, err := session.GetRoles(newAddr)
	conn.Commit(true)

	_, err = session.GrantRole(newAddr, big.NewInt(1))
	assert.NoError(t, err)
	conn.Commit(true)

	roles, err = session.GetRoles(newAddr)
	conn.Commit(true)

	fmt.Print(roles)
	assert.Nil(t, err)
}

func newMockEthConnection(gasLimit uint64, addressBalance *big.Int) (*c4tBackends.SimulatedBackend, common.Address, *c4tBind.TransactOpts, *ecdsa.PrivateKey) {
	adminKey, _ := crypto.GenerateKey()
	adminAuth, _ := c4tBind.NewKeyedTransactorWithChainID(adminKey, big.NewInt(1337))
	adminPublicKey := adminKey.Public()
	publicKeyECDSA, _ := adminPublicKey.(*ecdsa.PublicKey)
	adminAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	alloc := make(c4tCore.GenesisAlloc)

	// Add balance to admin
	alloc[adminAuth.From] = c4tCore.GenesisAccount{
		Balance: addressBalance,
	}

	return c4tBackends.NewSimulatedBackendWithInitialAdmin(rawdb.NewMemoryDatabase(), alloc, gasLimit, adminAddress), adminAddress, adminAuth, adminKey
}
