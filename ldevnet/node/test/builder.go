package test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-storedcounter"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/test"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/gen"
	genesis2 "github.com/filecoin-project/lotus/chain/gen/genesis"
	"github.com/filecoin-project/lotus/chain/messagepool"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	sectorstorage "github.com/filecoin-project/lotus/extern/sector-storage"
	"github.com/filecoin-project/lotus/extern/sector-storage/ffiwrapper"
	"github.com/filecoin-project/lotus/extern/sector-storage/mock"
	"github.com/filecoin-project/lotus/genesis"
	lotusminer "github.com/filecoin-project/lotus/miner"
	"github.com/filecoin-project/lotus/node"
	"github.com/filecoin-project/lotus/node/modules"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	testing2 "github.com/filecoin-project/lotus/node/modules/testing"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/storage/mockstorage"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	miner2 "github.com/filecoin-project/specs-actors/v2/actors/builtin/miner"
	"github.com/ipfs/go-datastore"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/multiformats/go-multiaddr"

	test_internal "github.com/textileio/lotus-devnet/ldevnet/api/test"
)

func init() {
	messagepool.HeadChangeCoalesceMinDelay = time.Microsecond
	messagepool.HeadChangeCoalesceMaxDelay = 2 * time.Microsecond
	messagepool.HeadChangeCoalesceMergeInterval = 100 * time.Nanosecond
}

func MockSbBuilder(fullOpts []test.FullNodeOpts, storage []test.StorageMiner, bigSector bool, ipfsAddr string, onlineMode bool) ([]test.TestNode, []test.TestStorageNode, error) {
	return mockSbBuilderOpts(fullOpts, storage, false, bigSector, ipfsAddr, onlineMode)
}

func mockSbBuilderOpts(fullOpts []test.FullNodeOpts, storage []test.StorageMiner, rpc bool, bigSector bool, ipfsAddr string, onlineMode bool) ([]test.TestNode, []test.TestStorageNode, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mn := mocknet.New(ctx)

	fulls := make([]test.TestNode, len(fullOpts))
	storers := make([]test.TestStorageNode, len(storage))

	var genbuf bytes.Buffer

	// PRESEAL SECTION, TRY TO REPLACE WITH BETTER IN THE FUTURE
	// TODO: would be great if there was a better way to fake the preseals
	var genms []genesis.Miner
	var genaccs []genesis.Actor
	var maddrs []address.Address
	var keys []*wallet.Key
	var pidKeys []crypto.PrivKey
	for i := 0; i < len(storage); i++ {
		maddr, err := address.NewIDAddress(genesis2.MinerStart + uint64(i))
		if err != nil {
			return nil, nil, err
		}

		preseals := storage[i].Preseal
		if preseals == test.PresealGenesis {
			preseals = test.GenesisPreseals
		}

		// <custom>
		size := abi.RegisteredSealProof_StackedDrg2KiBV1
		if bigSector {
			size = abi.RegisteredSealProof_StackedDrg512MiBV1
		}
		// </custom>

		genm, k, err := mockstorage.PreSeal(size, maddr, preseals)
		if err != nil {
			return nil, nil, err
		}

		pk, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		minerPid, err := peer.IDFromPrivateKey(pk)
		if err != nil {
			return nil, nil, err
		}

		genm.PeerId = minerPid

		wk, err := wallet.NewKey(*k)
		if err != nil {
			return nil, nil, err
		}

		genaccs = append(genaccs, genesis.Actor{
			Type:    genesis.TAccount,
			Balance: big.Mul(big.NewInt(400000000), types.NewInt(build.FilecoinPrecision)),
			Meta:    (&genesis.AccountMeta{Owner: wk.Address}).ActorMeta(),
		})

		keys = append(keys, wk)
		pidKeys = append(pidKeys, pk)
		maddrs = append(maddrs, maddr)
		genms = append(genms, *genm)
	}
	templ := &genesis.Template{
		Accounts:         genaccs,
		Miners:           genms,
		NetworkName:      "test",
		Timestamp:        uint64(time.Now().Add(-time.Hour * 100000).Unix()),
		VerifregRootKey:  gen.DefaultVerifregRootkeyActor,
		RemainderAccount: gen.DefaultRemainderAccountActor,
	}

	// END PRESEAL SECTION

	for i := 0; i < len(fullOpts); i++ {
		var genesis node.Option
		if i == 0 {
			genesis = node.Override(new(modules.Genesis), testing2.MakeGenesisMem(&genbuf, *templ))
		} else {
			genesis = node.Override(new(modules.Genesis), modules.LoadGenesis(genbuf.Bytes()))
		}

		_, err := node.New(ctx,
			node.FullAPI(&fulls[i].FullNode, node.Lite(fullOpts[i].Lite)),
			node.Online(),
			node.Repo(repo.NewMemory(nil)),
			node.MockHost(mn),
			node.Test(),

			node.Override(new(ffiwrapper.Verifier), mock.MockVerifier),

			// so that we subscribe to pubsub topics immediately
			node.Override(new(dtypes.Bootstrapper), dtypes.Bootstrapper(true)),

			genesis,

			fullOpts[i].Opts(fulls),

			node.ApplyIf(func(s *node.Settings) bool { return len(ipfsAddr) > 0 },
				node.Override(new(dtypes.ClientBlockstore), modules.IpfsClientBlockstore(ipfsAddr, onlineMode)),
				node.Override(new(dtypes.ClientRetrievalStoreManager), modules.ClientBlockstoreRetrievalStoreManager),
			),
		)
		if err != nil {
			return nil, nil, err
		}

		if rpc {
			fulls[i], err = fullRpc(fulls[i])
			if err != nil {
				return nil, nil, err
			}
		}
	}

	for i, def := range storage {
		// TODO: support non-bootstrap miners

		minerID := abi.ActorID(genesis2.MinerStart + uint64(i))

		if def.Full != 0 {
			return nil, nil, fmt.Errorf("storage nodes only supported on the first full node")
		}

		f := fulls[def.Full]
		if _, err := f.FullNode.WalletImport(ctx, &keys[i].KeyInfo); err != nil {
			return nil, nil, err
		}
		if err := f.FullNode.WalletSetDefault(ctx, keys[i].Address); err != nil {
			return nil, nil, err
		}

		sectors := make([]abi.SectorID, len(genms[i].Sectors))
		for i, sector := range genms[i].Sectors {
			sectors[i] = abi.SectorID{
				Miner:  minerID,
				Number: sector.SectorID,
			}
		}

		storers[i] = CreateTestStorageNode(ctx, genms[i].Worker, maddrs[i], pidKeys[i], f, mn, node.Options(
			node.Override(new(sectorstorage.SectorManager), func() (sectorstorage.SectorManager, error) {
				return mock.NewMockSectorMgr(sectors), nil
			}),
			node.Override(new(ffiwrapper.Verifier), mock.MockVerifier),
			node.Unset(new(*sectorstorage.Manager)),
		))

		if rpc {
			var err error
			storers[i], err = storerRpc(storers[i])
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if err := mn.LinkAll(); err != nil {
		return nil, nil, err
	}

	if len(storers) > 0 {
		// Mine 2 blocks to setup some CE stuff in some actors
		var wait sync.Mutex
		wait.Lock()

		test_internal.MineUntilBlock(ctx, fulls[0], storers[0], func(abi.ChainEpoch) {
			wait.Unlock()
		})
		wait.Lock()
		test_internal.MineUntilBlock(ctx, fulls[0], storers[0], func(abi.ChainEpoch) {
			wait.Unlock()
		})
		wait.Lock()
	}

	return fulls, storers, nil
}

func CreateTestStorageNode(ctx context.Context, waddr address.Address, act address.Address, pk crypto.PrivKey, tnd test.TestNode, mn mocknet.Mocknet, opts node.Option) test.TestStorageNode {
	r := repo.NewMemory(nil)

	lr, err := r.Lock(repo.StorageMiner)
	if err != nil {
		panic(err)
	}

	ks, err := lr.KeyStore()
	if err != nil {
		panic(err)
	}

	kbytes, err := pk.Bytes()
	if err != nil {
		panic(err)
	}

	err = ks.Put("libp2p-host", types.KeyInfo{
		Type:       "libp2p-host",
		PrivateKey: kbytes,
	})
	if err != nil {
		panic(err)
	}

	ds, err := lr.Datastore(ctx, "/metadata")
	if err != nil {
		panic(err)
	}
	err = ds.Put(datastore.NewKey("miner-address"), act.Bytes())
	if err != nil {
		panic(err)
	}

	nic := storedcounter.New(ds, datastore.NewKey(modules.StorageCounterDSPrefix))
	for i := 0; i < test.GenesisPreseals; i++ {
		_, err := nic.Next()
		if err != nil {
			panic(err)
		}
	}
	_, err = nic.Next()
	if err != nil {
		panic(err)
	}

	err = lr.Close()
	if err != nil {
		panic(err)
	}

	peerid, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		panic(err)
	}

	enc, err := actors.SerializeParams(&miner2.ChangePeerIDParams{NewID: abi.PeerID(peerid)})
	if err != nil {
		panic(err)
	}

	msg := &types.Message{
		To:     act,
		From:   waddr,
		Method: builtin.MethodsMiner.ChangePeerID,
		Params: enc,
		Value:  types.NewInt(0),
	}

	_, err = tnd.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		panic(err)
	}

	// start node
	var minerapi api.StorageMiner

	mineBlock := make(chan lotusminer.MineReq)
	_, err = node.New(ctx,
		node.StorageMiner(&minerapi),
		node.Online(),
		node.Repo(r),
		node.Test(),

		node.MockHost(mn),

		node.Override(new(api.FullNode), tnd),
		node.Override(new(*lotusminer.Miner), lotusminer.NewTestMiner(mineBlock, act)),

		opts,
	)
	if err != nil {
		panic(err)
	}

	/*// Bootstrap with full node
	remoteAddrs, err := tnd.NetAddrsListen(ctx)
	require.NoError(t, err)

	err = minerapi.NetConnect(ctx, remoteAddrs)
	require.NoError(t, err)*/
	mineOne := func(ctx context.Context, req lotusminer.MineReq) error {
		select {
		case mineBlock <- req:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return test.TestStorageNode{StorageMiner: minerapi, MineOne: mineOne}
}

func storerRpc(nd test.TestStorageNode) (test.TestStorageNode, error) {
	ma, listenAddr, err := CreateRPCServer(nd)
	if err != nil {
		return test.TestStorageNode{}, err
	}

	var storer test.TestStorageNode
	storer.StorageMiner, _, err = client.NewStorageMinerRPC(context.Background(), listenAddr, nil)
	if err != nil {
		return test.TestStorageNode{}, err
	}

	storer.ListenAddr = ma
	storer.MineOne = nd.MineOne
	return storer, nil
}
func fullRpc(nd test.TestNode) (test.TestNode, error) {
	ma, listenAddr, err := CreateRPCServer(nd)
	if err != nil {
		return test.TestNode{}, err
	}

	var full test.TestNode
	full.FullNode, _, err = client.NewFullNodeRPC(context.Background(), listenAddr, nil)
	if err != nil {
		return test.TestNode{}, err
	}

	full.ListenAddr = ma
	return full, nil
}

func CreateRPCServer(handler interface{}) (multiaddr.Multiaddr, string, error) {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("Filecoin", handler)
	testServ := httptest.NewServer(rpcServer) //  todo: close

	addr := testServ.Listener.Addr()
	listenAddr := "ws://" + addr.String()
	ma, err := parseWSMultiAddr(addr)
	if err != nil {
		return nil, "", err
	}
	return ma, listenAddr, err
}

func parseWSMultiAddr(addr net.Addr) (multiaddr.Multiaddr, error) {
	host, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil, err
	}
	ma, err := multiaddr.NewMultiaddr("/ip4/" + host + "/" + addr.Network() + "/" + port + "/ws")
	if err != nil {
		return nil, err
	}
	return ma, nil
}
