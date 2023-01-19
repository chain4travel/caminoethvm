package evm

import (
	"math/big"
	"testing"

	gconstants "github.com/ava-labs/coreth/constants"

	"github.com/ava-labs/avalanchego/utils/units"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	caminoGenesis = "{\"config\":{\"chainId\":502},\"initialAdmin\":\"0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC\", \"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"0100000000000000000000000000000000000000\":{\"code\":\"0x7300000000000000000000000000000000000000003014608060405260043610603d5760003560e01c80631e010439146042578063b6510bb314606e575b600080fd5b605c60048036036020811015605657600080fd5b503560b1565b60408051918252519081900360200190f35b818015607957600080fd5b5060af60048036036080811015608e57600080fd5b506001600160a01b03813516906020810135906040810135906060013560b6565b005b30cd90565b836001600160a01b031681836108fc8690811502906040516000604051808303818888878c8acf9550505050505015801560f4573d6000803e3d6000fd5b505050505056fea26469706673582212201eebce970fe3f5cb96bf8ac6ba5f5c133fc2908ae3dcd51082cfee8f583429d064736f6c634300060a0033\",\"balance\":\"0x0\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
)

func TestEVMStateTransfer(t *testing.T) {
	type args struct {
		balance     int64
		slotBalance int64
		amount      uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy path",
			args: args{
				balance:     10 * int64(units.MegaAvax),
				slotBalance: 1 * int64(units.MegaAvax),
				amount:      10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, vm, _, _, _ := GenesisVM(t, true, caminoGenesis, "", "")
			state, err := vm.blockChain.State()
			require.NoError(t, err)

			// Add balance to coinbase address
			state.AddBalance(gconstants.BlackholeAddr, big.NewInt(tt.args.balance))
			state.SetState(gconstants.BlackholeAddr, BalanceSlot, common.BigToHash(big.NewInt(tt.args.slotBalance)))

			// Cal the rewards tx
			tx, err := vm.NewCollectRewardsTx(tt.args.amount)
			require.NoError(t, err)

			err = tx.EVMStateTransfer(vm.ctx, state)
			require.NoError(t, err)
		})
	}
}
