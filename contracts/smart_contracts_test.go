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

	"github.com/chain4travel/caminoethvm/params"
)

var (
	ctx = context.Background()

	gasLimit = uint64(1)

	errAccessDeniedMsg = "execution reverted: Access denied"
	errUnkownRoleMsg   = "execution reverted: Unknown Role"
	errNotContractMsg  = "execution reverted: Not a contract"

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

	ADMIN_ROLE     = big.NewInt(1)
	GAS_FEE_ROLE   = big.NewInt(2)
	KYC_ROLE       = big.NewInt(4)
	BLACKLIST_ROLE = big.NewInt(8)
)

func TestAdminRoleFunctions(t *testing.T) {
	contractAddr := AdminProxyAddr

	// Initialize FilterOpts for event filtering
	filterOpts := bind.FilterOpts{Context: ctx, Start: 1, End: nil}

	// Initialize TransactOpts for each key
	adminOpts, err := bind.NewKeyedTransactorWithChainID(adminKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(alloc, gasLimit, adminAddr)
	defer func() {
		err := sim.Close()
		assert.NoError(t, err)
	}()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	adminSession := admin.BuildSession{Contract: adminContract, TransactOpts: *adminOpts}

	adminRole, err := adminSession.HasRole(adminAddr, ADMIN_ROLE)
	assert.NoError(t, err)

	// Assertion to check if Initial Admin address has indeed the admin role
	assert.True(t, adminRole)

	// Initialize role in every address (test GrantRole function)
	expectedRoles := addAndVerifyRoles(t, adminSession, sim)

	// Test Revoke Role
	// Delete blacklist role from blacklist address
	_, err = adminSession.RevokeRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	// Validate that there is no role in the address
	role, err := adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(role))))

	// Add the blacklist role again
	_, err = adminSession.GrantRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err = adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, BLACKLIST_ROLE, role)
	expectedRoles = append(expectedRoles, role)

	// Test Upgrade function
	// Initially, with an non-contract address. Should fail.
	_, err = adminSession.Upgrade(adminAddr)
	assert.EqualError(t, err, errNotContractMsg)

	// Then, with a contract address
	_, err = adminSession.Upgrade(contractAddr)
	assert.NoError(t, err)

	// Event Testing
	dropCounter := 0 // Counter of the total drop events
	dropItr, err := adminContract.FilterDropRole(&filterOpts)
	assert.NoError(t, err)
	assert.NotNil(t, dropItr)

	for dropItr.Next() {
		dropCounter++
		event := dropItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.Role, BLACKLIST_ROLE)
	}
	// Only one drop event has been emitted (Revoke Role) therefore, counter's value should be 1
	assert.EqualValues(t, 1, dropCounter)

	setCounter := 0 // Counter of the total set role events
	setItr, err := adminContract.FilterSetRole(&filterOpts)
	assert.NoError(t, err)
	assert.NotNil(t, setItr)

	for setItr.Next() {
		event := setItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.Role, expectedRoles[setCounter])
		setCounter++
	}
	// Four Set Roles events have been emitted in total
	// (including the events in addAndVerifyRoles function) therefore, counter's value should be 4
	assert.EqualValues(t, 4, setCounter)
}

func TestKycRoleFunctions(t *testing.T) {
	contractAddr := AdminProxyAddr

	// Initialize FilterOpts for event filtering
	filterOpts := bind.FilterOpts{Context: ctx, Start: 1, End: nil}

	// Initialize TransactOpts for each key
	kycOpts, err := bind.NewKeyedTransactorWithChainID(kycKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(alloc, gasLimit, kycAddr)
	defer func() {
		err := sim.Close()
		assert.NoError(t, err)
	}()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	kycSession := admin.BuildSession{Contract: adminContract, TransactOpts: *kycOpts}

	// Add KYC role in the address
	_, err = kycSession.GrantRole(kycAddr, KYC_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err := kycSession.GetRoles(kycAddr)
	assert.NoError(t, err)
	// At this point, kycAddr has both Admin and KYC Role
	assert.Equal(t, big.NewInt(5), role)

	st, err := kycSession.GetKycState(kycAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(st))))

	_, err = kycSession.ApplyKycState(kycAddr, false, big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	st, err = kycSession.GetKycState(kycAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), st)

	// Event Testing
	setCounter := 0 // Counter of the total set role events
	setItr, err := adminContract.FilterSetRole(&filterOpts)
	assert.NoError(t, err)
	assert.NotNil(t, setItr)

	for setItr.Next() {
		setCounter++
		event := setItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.Role, KYC_ROLE)
	}
	// Only one set role event has been emitted therefore, counter's value should be 1
	assert.EqualValues(t, 1, setCounter)

	kycCounter := 0 // Counter of the total set role events
	kycAddresses := make([]common.Address, 1)
	kycAddresses = append(kycAddresses, kycAddr)
	kycItr, err := adminContract.FilterKycStateChanged(&filterOpts, kycAddresses)
	assert.NoError(t, err)
	assert.NotNil(t, kycItr)

	for kycItr.Next() {
		kycCounter++
		event := kycItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.NewState, big.NewInt(1))
	}
	// Only one KYC State Changed event has been emitted therefore, counter's value should be 1
	assert.EqualValues(t, 1, kycCounter)
}

func TestGasFeeFunctions(t *testing.T) {
	contractAddr := AdminProxyAddr

	// Initialize FilterOpts for event filtering
	filterOpts := bind.FilterOpts{Context: ctx, Start: 1, End: nil}

	// Initialize TransactOpts for each key
	gasFeeOpts, err := bind.NewKeyedTransactorWithChainID(gasFeeKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(alloc, gasLimit, gasFeeAddr)
	defer func() {
		err := sim.Close()
		assert.NoError(t, err)
	}()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	gasFeeSession := admin.BuildSession{Contract: adminContract, TransactOpts: *gasFeeOpts}

	// Add Gas Fee Role in the address
	_, err = gasFeeSession.GrantRole(gasFeeAddr, GAS_FEE_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err := gasFeeSession.GetRoles(gasFeeAddr)
	assert.NoError(t, err)
	// At this point, gasFeeAddr has both Admin and Gas Fee Role
	assert.Equal(t, big.NewInt(3), role)

	bf, err := gasFeeSession.GetBaseFee()
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(bf))))

	_, err = gasFeeSession.SetBaseFee(big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	bf, err = gasFeeSession.GetBaseFee()
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), bf)

	// Event Testing
	setCounter := 0 // Counter of the total set role events
	setItr, err := adminContract.FilterSetRole(&filterOpts)
	assert.NoError(t, err)
	assert.NotNil(t, setItr)

	for setItr.Next() {
		setCounter++
		event := setItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.Role, GAS_FEE_ROLE)
	}
	// Only one set role event has been emitted therefore, counter's value should be 1
	assert.EqualValues(t, 1, setCounter)

	gasCounter := 0 // Counter of the total set role events
	gasItr, err := adminContract.FilterGasFeeSet(&filterOpts)
	assert.NoError(t, err)
	assert.NotNil(t, gasItr)

	for gasItr.Next() {
		gasCounter++
		event := gasItr.Event
		assert.NotNil(t, event)
		assert.EqualValues(t, event.NewGasFee, big.NewInt(1))
	}
	// Only one KYC State Changed event has been emitted therefore, counter's value should be 1
	assert.EqualValues(t, 1, gasCounter)
}

func TestBlacklistFunctions(t *testing.T) {
	var signature [4]byte
	contractAddr := AdminProxyAddr

	// Initialize TransactOpts for each key
	blacklistOpts, err := bind.NewKeyedTransactorWithChainID(blacklistKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(alloc, gasLimit, blacklistAddr)
	defer func() {
		err := sim.Close()
		assert.NoError(t, err)
	}()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	blacklistSession := admin.BuildSession{Contract: adminContract, TransactOpts: *blacklistOpts}

	// Add Gas Fee Role in the address
	_, err = blacklistSession.GrantRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	sim.Commit(true)

	role, err := blacklistSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	// At this point, blacklistAddr has both Admin and Blacklist Role
	assert.Equal(t, big.NewInt(9), role)

	blst, err := blacklistSession.GetBlacklistState(blacklistAddr, signature)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), big.NewInt(int64(common.Big0.Cmp(blst))))

	_, err = blacklistSession.SetBlacklistState(blacklistAddr, signature, big.NewInt(1))
	assert.NoError(t, err)

	sim.Commit(true)

	blst, err = blacklistSession.GetBlacklistState(blacklistAddr, signature)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(1), blst)
}

func TestDummySession(t *testing.T) {
	contractAddr := AdminProxyAddr

	// Initialize FilterOpts for event filtering
	filterOpts := bind.FilterOpts{Context: ctx, Start: 1, End: nil}

	// Initialize TransactOpts for each key
	dummyOpts, err := bind.NewKeyedTransactorWithChainID(dummyKey, big.NewInt(1337))
	assert.NoError(t, err)

	// Generate GenesisAlloc
	alloc := makeGenesisAllocation()

	// Generate SimulatedBackend
	sim := backends.NewSimulatedBackendWithInitialAdmin(alloc, gasLimit, adminAddr)
	defer func() {
		err := sim.Close()
		assert.NoError(t, err)
	}()

	sim.Commit(true)

	adminContract, err := admin.NewBuild(contractAddr, sim)
	assert.NoError(t, err)

	// BuildSession Initialization
	dummySession := admin.BuildSession{Contract: adminContract, TransactOpts: *dummyOpts}

	// At this point, dummyAddr has no roles and therefore, should not be able to call any function
	_, err = dummySession.SetBaseFee(big.NewInt(1))
	assert.EqualError(t, err, errAccessDeniedMsg)

	sim.Commit(true)

	_, err = dummySession.ApplyKycState(dummyAddr, false, big.NewInt(1))
	assert.EqualError(t, err, errAccessDeniedMsg)

	sim.Commit(true)

	_, err = dummySession.GrantRole(dummyAddr, ADMIN_ROLE)
	assert.EqualError(t, err, errAccessDeniedMsg)

	// Event Testing
	// No events have been emitted so normally, nothing should be filtered
	setItr, err := adminContract.FilterSetRole(&filterOpts)
	assert.NoError(t, err)
	assert.False(t, setItr.Next())

	dropItr, err := adminContract.FilterDropRole(&filterOpts)
	assert.NoError(t, err)
	assert.False(t, dropItr.Next())

	gasItr, err := adminContract.FilterGasFeeSet(&filterOpts)
	assert.NoError(t, err)
	assert.False(t, gasItr.Next())

	kycAddresses := make([]common.Address, 1)
	kycAddresses = append(kycAddresses, kycAddr)
	kycItr, err := adminContract.FilterKycStateChanged(&filterOpts, kycAddresses)
	assert.NoError(t, err)
	assert.False(t, kycItr.Next())
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

func addAndVerifyRoles(t *testing.T, adminSession admin.BuildSession, sim *backends.SimulatedBackend) []*big.Int {
	roles := make([]*big.Int, 0)
	// Add and verify Roles
	_, err := adminSession.GrantRole(kycAddr, KYC_ROLE)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(gasFeeAddr, GAS_FEE_ROLE)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(blacklistAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	_, err = adminSession.GrantRole(dummyAddr, big.NewInt(100))
	assert.EqualError(t, err, errUnkownRoleMsg)

	sim.Commit(true)

	kycRole, err := adminSession.GetRoles(kycAddr)
	assert.NoError(t, err)
	assert.Equal(t, KYC_ROLE, kycRole)
	roles = append(roles, kycRole)

	gasFeeRole, err := adminSession.GetRoles(gasFeeAddr)
	assert.NoError(t, err)
	assert.Equal(t, GAS_FEE_ROLE, gasFeeRole)
	roles = append(roles, gasFeeRole)

	blacklistRole, err := adminSession.GetRoles(blacklistAddr)
	assert.NoError(t, err)
	assert.Equal(t, BLACKLIST_ROLE, blacklistRole)
	roles = append(roles, blacklistRole)

	return roles
}
