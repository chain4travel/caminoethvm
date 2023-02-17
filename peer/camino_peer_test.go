package peer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	msg "github.com/ava-labs/avalanchego/message"
	"github.com/ava-labs/avalanchego/proto/pb/p2p"
	"github.com/ava-labs/coreth/plugin/evm/message"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []string{
		"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",
	}

	codec := buildCodec(t,
		message.BlockRequest{},
		message.BlockResponse{},
		message.CodeRequest{},
		message.CodeResponse{},
		message.LeafsRequest{},
		message.LeafsResponse{},
		message.SyncSummary{},
		message.AtomicTxGossip{},
		message.EthTxsGossip{},
		types.Block{},
	)

	mc, err := msg.NewCreator(
		prometheus.NewRegistry(),
		"",
		true,
		10*time.Second,
	)
	require.NoError(t, err, "Unable to create message creator")

	for i, msg := range tests {
		t.Run(fmt.Sprintf("Decoding-%d", i), func(t *testing.T) {
			data, err := base64.StdEncoding.DecodeString(msg)
			require.NoError(t, err)
			//fmt.Println(len(data))

			ibmsg, err := mc.Parse(data, ids.NodeID(ids.ShortEmpty), func() { fmt.Println("onFinishedHandling 👋") })
			require.NoError(t, err)
			smth := ibmsg.Message()

			appGosipMsg, ok := smth.(*p2p.AppGossip)
			require.True(t, ok)

			dec3 := message.BlockResponse{}
			_, err = codec.Unmarshal(appGosipMsg.GetAppBytes(), &dec3)
			require.NoError(t, err)

			blks := dec3.Blocks
			fmt.Println(len(blks))

			for _, bb := range blks {
				blk := types.Block{}
				stream := rlp.NewStream(bytes.NewReader(bb), 0)
				err = blk.DecodeRLP(stream)
				require.NoError(t, err)
			}
		})
	}
}
