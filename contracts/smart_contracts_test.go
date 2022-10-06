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

package contracts

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/chain4travel/caminoethvm/accounts/abi/bind"
	"github.com/chain4travel/caminoethvm/accounts/abi/bind/backends"
	admin "github.com/chain4travel/caminoethvm/contracts/build_contracts/admin/src"
	"github.com/chain4travel/caminoethvm/core"

	"github.com/chain4travel/caminoethvm/core/rawdb"
	"github.com/chain4travel/caminoethvm/params"
)

var (
	ctx      = context.Background()
	gasLimit = uint64(1)

	errAccessDeniedMsg = "execution reverted: Access denied"

	adminKey, _     = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	kycKey, _       = crypto.HexToECDSA("0e6de2b744bd97ab6abd2e8fc624befbac7ed5b37e8a7e8ddd164e23d7ac06be")
	gasFeeKey, _    = crypto.HexToECDSA("04214cc61e1feaf005aa25b7771d33ca5c4aea959d21fe9a1429f822fa024171")
	blacklistKey, _ = crypto.HexToECDSA("b32d5aa5b8f4028218538c8c5b14b5c14f3f2e35b236e4bbbff09b669e69e46c")
	dummyKey, _     = crypto.HexToECDSA("62802c57c0e3c24ae0ce354f7d19f7659ddbe506547b00e9e6a722980d2fed3d")

	adminAddr     = crypto.PubkeyToAddress(adminKey.PublicKey)
	kycAddr       = crypto.PubkeyToAddress(kycKey.PublicKey)
	gasFeeAddr    = crypto.PubkeyToAddress(gasFeeKey.PublicKey)
	blacklistAddr = crypto.PubkeyToAddress(blacklistKey.PublicKey)
	dummyAddr     = crypto.PubkeyToAddress(dummyKey.PublicKey)

	AdminProxyAddr = common.HexToAddress("0x010000000000000000000000000000000000000a")

	NO_ROLE        = big.NewInt(0)
	ADMIN_ROLE     = big.NewInt(1)
	GAS_FEE_ROLE   = big.NewInt(2)
	KYC_ROLE       = big.NewInt(4)
	BLACKLIST_ROLE = big.NewInt(8)
)

func TestSmartContracts(t *testing.T) {
	contractAddr := AdminProxyAddr

	// Initialize TransactOpts for each key
	adminOpts, err := bind.NewKeyedTransactorWithChainID(adminKey, big.NewInt(1337))
	assert.NoError(t, err)

	kycOpts, err := bind.NewKeyedTransactorWithChainID(kycKey, big.NewInt(1337))
	assert.NoError(t, err)

	gasFeeOpts, err := bind.NewKeyedTransactorWithChainID(gasFeeKey, big.NewInt(1337))
	assert.NoError(t, err)

	dummyOpts, err := bind.NewKeyedTransactorWithChainID(dummyKey, big.NewInt(1337))
	assert.NoError(t, err)

	blacklistOpts, err := bind.NewKeyedTransactorWithChainID(blacklistKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(rawdb.NewMemoryDatabase(), alloc, gasLimit, adminAddr)
	defer sim.Close()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	adminSession := admin.BuildSession{Contract: adminContract, TransactOpts: *adminOpts}
	kycSession := admin.BuildSession{Contract: adminContract, TransactOpts: *kycOpts}
	gasFeeSession := admin.BuildSession{Contract: adminContract, TransactOpts: *gasFeeOpts}
	dummySession := admin.BuildSession{Contract: adminContract, TransactOpts: *dummyOpts}
	blacklistSession := admin.BuildSession{Contract: adminContract, TransactOpts: *blacklistOpts}

	adminRole, err := adminSession.HasRole(adminAddr, ADMIN_ROLE)
	assert.NoError(t, err)

	// Assertion to check if Initial Admin adress has indeed the admin role
	assert.True(t, adminRole)

	// Initialize role in every address
	addAndVerifyRoles(t, adminSession, sim)

	// Test Admin Only Functions
	testAdminRoleFunctions(t, adminSession, adminOpts, sim)

	// Test KYC Only Functions
	testKycRoleFunctions(t, kycSession, adminOpts, sim)

	// Test Gas Fee Only Functions
	testGasFeeFunctions(t, gasFeeSession, adminOpts, sim)

	// Test Blacklist Only Functions
	testBlacklistFunctions(t, blacklistSession, adminOpts, sim)

	// Test Functions with dummy address. Should fail
	testDummySession(t, dummySession, adminOpts, sim)
}

func addAndVerifyRoles(t *testing.T, adminSession admin.BuildSession, sim *backends.SimulatedBackend) {
	// Add and verify Roles
	_, err := adminSession.GrantRole(kycAddr, KYC_ROLE)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(gasFeeAddr, GAS_FEE_ROLE)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	gasFeeRole, err := adminSession.GetRoles(gasFeeAddr)
	assert.NoError(t, err)
	assert.Equal(t, GAS_FEE_ROLE, gasFeeRole)

	kycRole, err := adminSession.GetRoles(kycAddr)
	assert.NoError(t, err)
	assert.Equal(t, KYC_ROLE, kycRole)

	blacklistRole, err := adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	assert.Equal(t, BLACKLIST_ROLE, blacklistRole)
}

func testAdminRoleFunctions(t *testing.T, adminSession admin.BuildSession, opts *bind.TransactOpts, sim *backends.SimulatedBackend) {
	_, err := adminSession.RevokeRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err := adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err = adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, BLACKLIST_ROLE, role)
}

func testKycRoleFunctions(t *testing.T, kycSession admin.BuildSession, opts *bind.TransactOpts, sim *backends.SimulatedBackend) {
	role, err := kycSession.GetRoles(kycAddr)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(4), role)

	st, err := kycSession.GetKycState(kycAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(st))))

	_, err = kycSession.ApplyKycState(kycAddr, false, big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	st, err = kycSession.GetKycState(kycAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), st)
}

func testGasFeeFunctions(t *testing.T, gasFeeSession admin.BuildSession, opts *bind.TransactOpts, sim *backends.SimulatedBackend) {
	role, err := gasFeeSession.GetRoles(gasFeeAddr)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(2), role)

	bf, err := gasFeeSession.GetBaseFee()
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(bf))))

	_, err = gasFeeSession.SetBaseFee(big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	bf, err = gasFeeSession.GetBaseFee()
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), bf)
}

func testBlacklistFunctions(t *testing.T, blacklistSession admin.BuildSession, opts *bind.TransactOpts, sim *backends.SimulatedBackend) {
	var tmp [4]byte

	blst, err := blacklistSession.GetBlacklistState(blacklistAddr, tmp)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(blst))))

	_, err = blacklistSession.SetBlacklistState(blacklistAddr, tmp, big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	blst, err = blacklistSession.GetBlacklistState(blacklistAddr, tmp)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), blst)
}

func testDummySession(t *testing.T, dummySession admin.BuildSession, opts *bind.TransactOpts, sim *backends.SimulatedBackend) {
	var tmp [4]byte

	_, err := dummySession.SetBlacklistState(dummyAddr, tmp, big.NewInt(1))
	assert.EqualError(t, err, errAccessDeniedMsg)

	sim.Commit(true)

	_, err = dummySession.SetBaseFee(big.NewInt(1))
	assert.EqualError(t, err, errAccessDeniedMsg)

	sim.Commit(true)

	_, err = dummySession.ApplyKycState(dummyAddr, false, big.NewInt(1))
	assert.EqualError(t, err, errAccessDeniedMsg)

	sim.Commit(true)
}

func makeGenesisAllocation() core.GenesisAlloc {
	alloc := make(core.GenesisAlloc)

	alloc[adminAddr] = core.GenesisAccount{Balance: big.NewInt(params.Ether)}
	alloc[kycAddr] = core.GenesisAccount{Balance: big.NewInt(params.Ether)}
	alloc[gasFeeAddr] = core.GenesisAccount{Balance: big.NewInt(params.Ether)}
	alloc[blacklistAddr] = core.GenesisAccount{Balance: big.NewInt(params.Ether)}
	alloc[dummyAddr] = core.GenesisAccount{Balance: big.NewInt(params.Ether)}

	return alloc
}
