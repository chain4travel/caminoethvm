package admin_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/accounts/abi/bind/backends"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

type waitDeployTest struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}

func waitDeployTestExec(name string, test waitDeployTest, backend *backends.SimulatedBackend, t *testing.T) {
	// Create the transaction
	head, _ := backend.HeaderByNumber(context.Background(), nil) // Should be child's, good enough
	gasPrice := new(big.Int).Add(head.BaseFee, big.NewInt(1))

	tx := types.NewContractCreation(0, big.NewInt(0), test.gas, gasPrice, common.FromHex(test.code))
	signer := types.NewLondonSigner(big.NewInt(1337))
	tx, _ = types.SignTx(tx, signer, testKey)

	// Wait for it to get mined in the background.
	var (
		err     error
		address common.Address
		mined   = make(chan struct{})
		ctx     = context.Background()
	)
	go func() {
		address, err = bind.WaitDeployed(ctx, backend, tx)
		close(mined)
	}()

	// Send and mine the transaction.
	if err := backend.SendTransaction(ctx, tx); err != nil {
		t.Errorf("Failed to send transaction: %s", err)
	}
	backend.Commit(true)

	select {
	case <-mined:
		if err != test.wantErr {
			t.Errorf("test %q: error mismatch: want %q, got %q", name, test.wantErr, err)
		}
		if address != test.wantAddress {
			t.Errorf("test %q: unexpected contract address %s", name, address.Hex())
		}
	case <-time.After(2 * time.Second):
		t.Errorf("test %q: timeout", name)
	}
}

func TestWaitDeployedVerified(t *testing.T) {
	tests := map[string]waitDeployTest{
		"successful deploy": {
			code:        `6060604052600a8060106000396000f360606040526008565b00`,
			gas:         3000000,
			wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
		},
		"empty code": {
			code:        ``,
			gas:         300000,
			wantErr:     bind.ErrNoCodeAfterDeploy,
			wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
		},
	}
	for name, test := range tests {
		backend := backends.NewSimulatedBackendWithKYCVerified(
			core.GenesisAlloc{
				crypto.PubkeyToAddress(testKey.PublicKey): {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
			},
			10000000,
			crypto.PubkeyToAddress(testKey.PublicKey),
		)
		defer backend.Close()

		waitDeployTestExec(name, test, backend, t)
	}
}

func TestWaitDeployUnverified(t *testing.T) {
	tests := map[string]waitDeployTest{
		"failed deploy": {
			code:        `6060604052600a8060106000396000f360606040526008565b00`,
			gas:         3000000,
			wantErr:     bind.ErrNoCodeAfterDeploy,
			wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
		},
	}
	for name, test := range tests {
		backend := backends.NewSimulatedBackendWithEnforcementEnabled(
			core.GenesisAlloc{
				crypto.PubkeyToAddress(testKey.PublicKey): {Balance: new(big.Int).Mul(big.NewInt(10000000000000000), big.NewInt(1000))},
			},
			10000000,
		)
		defer backend.Close()

		waitDeployTestExec(name, test, backend, t)
	}
}
