package admin_test

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/accounts/abi/bind/backends"
	adminConcract "github.com/ava-labs/coreth/contracts/build_contracts/admin/src"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/admin"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	ctx        = context.Background()
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr   = crypto.PubkeyToAddress(testKey.PublicKey)

	BLACKLIST_ROLE = big.NewInt(8)

	/*
		pragma solidity ^0.8.0;

		contract Test {

			mapping(uint32 => uint32) private testMap;

			function foo() public pure { return; }
			function bar() public pure { foo(); }
			function testAdd(uint32 a, uint32 b) public returns (uint32) {
				return a+b;
			}
			function testSave(uint32 a) public {
				testMap[a] = a;
			}
			function pack(address account, bytes4 signature) public pure returns (uint256) {
				return uint256(uint160(account)) | (uint256(uint32(signature)) << 160);
			}
		}
	*/
	contractAbi = "[{\"inputs\":[],\"name\":\"bar\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"foo\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"signature\",\"type\":\"bytes4\"}],\"name\":\"pack\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"testSave\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	contractBin = "608060405234801561001057600080fd5b5061015f806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633105badb1461005157806396e181cf1461008d578063c2985578146100ac578063febb0f7e146100ae575b600080fd5b61007b61005f3660046100d7565b60401c63ffffffff60a01b166001600160a01b03919091161790565b60405190815260200160405180910390f35b6100ac336000908152602081905260409020805460ff19166001179055565b005b6100ac6100b8565b565b6100b6336000908152602081905260409020805460ff19166001179055565b600080604083850312156100ea57600080fd5b82356001600160a01b038116811461010157600080fd5b915060208301356001600160e01b03198116811461011e57600080fd5b80915050925092905056fea26469706673582212206e1b710b9bad5e16818e269598447f511b9750832ff98ea565104a8e518dbab864736f6c63430008110033"
)

func deployTestContract(t *testing.T, backend *backends.SimulatedBackend) (common.Address, error) {
	opts, _ := bind.NewKeyedTransactorWithChainID(testKey, big.NewInt(1337))

	parsed, _ := abi.JSON(strings.NewReader(contractAbi))
	addr, _, _, err := bind.DeployContract(opts, parsed, common.FromHex(contractBin), backend)
	return addr, err
}

func TestDeployUnverified(t *testing.T) {

	backend := backends.NewSimulatedBackendWithEnforcementEnabled(
		core.GenesisAlloc{
			testAddr: {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
		},
		10000000,
	)
	defer backend.Close()
	_, err := deployTestContract(t, backend)
	if err == nil {
		t.Fatalf("sucessfully created contract without kyc verification")
	}
}

func TestDeployVerified(t *testing.T) {

	backend := backends.NewSimulatedBackendWithKYCVerified(
		core.GenesisAlloc{
			testAddr: {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
		},
		10000000,
		testAddr,
	)
	defer backend.Close()
	_, err := deployTestContract(t, backend)
	assert.NoError(t, err, "failed to deploy contract")

}

func blacklistFunction(t *testing.T, backend *backends.SimulatedBackend, addr common.Address, selector []byte) {
	// Initialize TransactOpts for each key
	blacklistOpts, err := bind.NewKeyedTransactorWithChainID(testKey, big.NewInt(1337))
	assert.NoError(t, err)

	backend.Commit(true)

	adminContract, err := adminConcract.NewBuild(admin.AdminContractAddr, backend)
	assert.NoError(t, err)

	// BuildSession Initialization
	blacklistSession := adminConcract.BuildSession{Contract: adminContract, TransactOpts: *blacklistOpts}

	// Add Blacklist Role in the address
	_, err = blacklistSession.GrantRole(testAddr, BLACKLIST_ROLE)
	assert.NoError(t, err)

	backend.Commit(true)

	var selector2 [4]byte
	copy(selector2[:], selector)

	// Add Blacklist Role in the address
	_, err = blacklistSession.SetBlacklistState(addr, selector2, common.Big1)
	assert.NoError(t, err)
	backend.Commit(true)

	a, err := blacklistSession.GetBlacklistState(addr, selector2)
	assert.NoError(t, err)
	assert.True(t, a.Cmp(common.Big1) == 0)
}

func TestFucntionDirectExecution(t *testing.T) {

	for _, test := range []struct {
		functionToBlacklist string
		functionToCall      string
	}{
		{"foo()", "foo()"},
		{"bar()", "foo()"},
		{"bar()", "testSave()"},
	} {
		backend := backends.NewSimulatedBackendWithKYCVerified(
			core.GenesisAlloc{
				testAddr: {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
			},
			10000000,
			testAddr,
		)
		defer backend.Close()

		blacklistSelector := crypto.Keccak256([]byte(test.functionToBlacklist))[:4]
		callSelector := crypto.Keccak256([]byte(test.functionToCall))[:4]

		contractAddr, err := deployTestContract(t, backend)
		assert.NoError(t, err, "failed to deploy contract")

		blacklistFunction(t, backend, contractAddr, blacklistSelector)

		backend.Commit(true)

		head, _ := backend.HeaderByNumber(ctx, nil) // Should be child's, good enough
		gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))
		tx := types.NewTransaction(3, contractAddr, common.Big0, 30000, gasPrice, callSelector[:])
		signer := types.NewLondonSigner(big.NewInt(1337))
		tx, _ = types.SignTx(tx, signer, testKey)

		err = backend.SendTransaction(ctx, tx)
		assert.NoError(t, err)

		backend.Commit(true)

		reciept, err := backend.TransactionReceipt(ctx, tx.Hash())
		assert.NoError(t, err)
		assert.Equal(t, reciept.Status, uint64(0))

	}

}
