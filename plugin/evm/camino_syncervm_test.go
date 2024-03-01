// Copyright (C) 2023, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	statesyncclient "github.com/ava-labs/coreth/sync/client"
	"github.com/ava-labs/coreth/sync/statesync"
	"github.com/ava-labs/coreth/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestSkipStateSyncCamino(t *testing.T) {
	rand.Seed(1)
	test := syncTest{
		syncableInterval:   256,
		stateSyncMinBlocks: 300, // must be greater than [syncableInterval] to skip sync
		syncMode:           block.StateSyncSkipped,
	}
	vmSetup := createSyncServerAndClientVMsCamino(t, test)

	testSyncerVM(t, vmSetup, test)
}

func createSyncServerAndClientVMsCamino(t *testing.T, test syncTest) *syncVMSetup {
	var (
		require      = require.New(t)
		importAmount = 2000000 * units.Avax // 2M avax
		alloc        = map[ids.ShortID]uint64{
			testShortIDAddrs[0]: importAmount,
		}
	)
	_, serverVM, _, serverAtomicMemory, serverAppSender := GenesisVMWithUTXOs(
		t, true, genesisJSONSunrisePhase0, "", "", alloc,
	)
	t.Cleanup(func() {
		log.Info("Shutting down server VM")
		require.NoError(serverVM.Shutdown(context.Background()))
	})

	var (
		importTx, exportTx *Tx
		err                error
	)
	generateAndAcceptBlocks(t, serverVM, parentsToGet, func(i int, gen *core.BlockGen) {
		switch i {
		case 0:
			// spend the UTXOs from shared memory
			importTx, err = serverVM.newImportTx(serverVM.ctx.XChainID, testEthAddrs[0], initialBaseFee, []*secp256k1.PrivateKey{testKeys[0]})
			require.NoError(err)
			require.NoError(serverVM.issueTx(importTx, true /*=local*/))
		case 1:
			// export some of the imported UTXOs to test exportTx is properly synced
			exportTx, err = serverVM.newExportTx(
				serverVM.ctx.AVAXAssetID,
				importAmount/2,
				serverVM.ctx.XChainID,
				testShortIDAddrs[0],
				initialBaseFee,
				[]*secp256k1.PrivateKey{testKeys[0]},
			)
			require.NoError(err)
			require.NoError(serverVM.issueTx(exportTx, true /*=local*/))
		default: // Generate simple transfer transactions.
			pk := testKeys[0].ToECDSA()
			tx := types.NewTransaction(gen.TxNonce(testEthAddrs[0]), testEthAddrs[1], common.Big1, params.TxGas, initialBaseFee, nil)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(serverVM.chainID), pk)
			require.NoError(err)
			gen.AddTx(signedTx)
		}
	})

	// override serverAtomicTrie's commitInterval so the call to [serverAtomicTrie.Index]
	// creates a commit at the height [syncableInterval]. This is necessary to support
	// fetching a state summary.
	serverAtomicTrie := serverVM.atomicTrie.(*atomicTrie)
	serverAtomicTrie.commitInterval = test.syncableInterval
	require.NoError(serverAtomicTrie.commit(test.syncableInterval, serverAtomicTrie.LastAcceptedRoot()))
	require.NoError(serverVM.db.Commit())

	serverSharedMemories := newSharedMemories(serverAtomicMemory, serverVM.ctx.ChainID, serverVM.ctx.XChainID)
	serverSharedMemories.assertOpsApplied(t, importTx.mustAtomicOps())
	serverSharedMemories.assertOpsApplied(t, exportTx.mustAtomicOps())

	// make some accounts
	trieDB := trie.NewDatabase(serverVM.chaindb)
	root, accounts := statesync.FillAccountsWithOverlappingStorage(t, trieDB, types.EmptyRootHash, 1000, 16)

	// patch serverVM's lastAcceptedBlock to have the new root
	// and update the vm's state so the trie with accounts will
	// be returned by StateSyncGetLastSummary
	lastAccepted := serverVM.blockChain.LastAcceptedBlock()
	patchedBlock := patchBlock(lastAccepted, root, serverVM.chaindb)
	blockBytes, err := rlp.EncodeToBytes(patchedBlock)
	require.NoError(err)
	internalBlock, err := serverVM.parseBlock(context.Background(), blockBytes)
	require.NoError(err)
	internalBlock.(*Block).SetStatus(choices.Accepted)
	require.NoError(serverVM.State.SetLastAcceptedBlock(internalBlock))

	// patch syncableInterval for test
	serverVM.StateSyncServer.(*stateSyncServer).syncableInterval = test.syncableInterval

	// initialise [syncerVM] with blank genesis state
	stateSyncEnabledJSON := fmt.Sprintf(`{"state-sync-enabled":true, "state-sync-min-blocks": %d}`, test.stateSyncMinBlocks)
	syncerEngineChan, syncerVM, syncerDBManager, syncerAtomicMemory, syncerAppSender := GenesisVMWithUTXOs(
		t, false, "", stateSyncEnabledJSON, "", alloc,
	)
	shutdownOnceSyncerVM := &shutdownOnceVM{VM: syncerVM}
	t.Cleanup(func() {
		require.NoError(shutdownOnceSyncerVM.Shutdown(context.Background()))
	})
	require.NoError(syncerVM.SetState(context.Background(), snow.StateSyncing))
	enabled, err := syncerVM.StateSyncEnabled(context.Background())
	require.NoError(err)
	require.True(enabled)

	// override [syncerVM]'s commit interval so the atomic trie works correctly.
	syncerVM.atomicTrie.(*atomicTrie).commitInterval = test.syncableInterval

	// override [serverVM]'s SendAppResponse function to trigger AppResponse on [syncerVM]
	serverAppSender.SendAppResponseF = func(ctx context.Context, nodeID ids.NodeID, requestID uint32, response []byte) error {
		if test.responseIntercept == nil {
			go syncerVM.AppResponse(ctx, nodeID, requestID, response)
		} else {
			go test.responseIntercept(syncerVM, nodeID, requestID, response)
		}

		return nil
	}

	// connect peer to [syncerVM]
	require.NoError(
		syncerVM.Connected(
			context.Background(),
			serverVM.ctx.NodeID,
			statesyncclient.StateSyncVersion,
		),
	)

	// override [syncerVM]'s SendAppRequest function to trigger AppRequest on [serverVM]
	syncerAppSender.SendAppRequestF = func(ctx context.Context, nodeSet set.Set[ids.NodeID], requestID uint32, request []byte) error {
		nodeID, hasItem := nodeSet.Pop()
		require.True(hasItem, "expected nodeSet to contain at least 1 nodeID")
		err := serverVM.AppRequest(ctx, nodeID, requestID, time.Now().Add(1*time.Second), request)
		require.NoError(err)
		return nil
	}

	return &syncVMSetup{
		serverVM:        serverVM,
		serverAppSender: serverAppSender,
		includedAtomicTxs: []*Tx{
			importTx,
			exportTx,
		},
		fundedAccounts:       accounts,
		syncerVM:             syncerVM,
		syncerDBManager:      syncerDBManager,
		syncerEngineChan:     syncerEngineChan,
		syncerAtomicMemory:   syncerAtomicMemory,
		shutdownOnceSyncerVM: shutdownOnceSyncerVM,
	}
}
