package ldevnet

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storedcounter"
	"github.com/filecoin-project/go-sectorbuilder"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/test"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors"
	genesis2 "github.com/filecoin-project/lotus/chain/gen/genesis"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/genesis"
	"github.com/filecoin-project/lotus/lib/jsonrpc"
	"github.com/filecoin-project/lotus/miner"
	"github.com/filecoin-project/lotus/node"
	"github.com/filecoin-project/lotus/node/modules"
	modtest "github.com/filecoin-project/lotus/node/modules/testing"
	"github.com/filecoin-project/lotus/node/repo"
	"github.com/filecoin-project/lotus/storage/sbmock"
	"github.com/filecoin-project/lotus/storage/sealmgr"
	"github.com/filecoin-project/lotus/storage/sealmgr/advmgr"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	saminer "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

var DefaultDuration = time.Millisecond * 100

func init() {
	build.SectorSizes = []abi.SectorSize{2048}
	power.ConsensusMinerMinPower = big.NewInt(2048)
	os.Setenv("TRUST_PARAMS", "1")
}

type LocalDevnet struct {
	Client *apistruct.FullNodeStruct

	numMiners int
	cancel    context.CancelFunc
	closer    func()
	done      chan struct{}
}

func (ld *LocalDevnet) Close() {
	ld.cancel()

	for i := 0; i < ld.numMiners; i++ {
		<-ld.done
	}
	close(ld.done)
	ld.closer()
}

func New(numMiners int, blockDur time.Duration) (*LocalDevnet, error) {
	miners := make([]int, numMiners)
	n, sn, closer, err := rpcBuilder(1, miners)
	if err != nil {
		return nil, err
	}
	client := n[0].FullNode.(*apistruct.FullNodeStruct)
	ctx, cancel := context.WithCancel(context.Background())
	addrinfo, err := client.NetAddrsListen(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	done := make(chan struct{})
	for i := range miners {
		if err := sn[i].NetConnect(ctx, addrinfo); err != nil {
			cancel()
			return nil, err
		}
		time.Sleep(time.Second)

		mine := i == 0
		go func(i int) {
			defer func() { done <- struct{}{} }()
			for mine {
				time.Sleep(blockDur)
				if ctx.Err() != nil {
					mine = false
					continue
				}
				if err := sn[i].MineOne(context.Background()); err != nil {
					panic(err)
				}
			}
		}(i)
	}
	for i := range miners {
		for j := range miners {
			if j == i {
				continue
			}
			mainfo, err := sn[j].NetAddrsListen(ctx)
			if err != nil {
				cancel()
				return nil, err
			}
			if err := sn[i].NetConnect(ctx, mainfo); err != nil {
				cancel()
				return nil, err
			}
		}
	}

	time.Sleep(blockDur * 5) // Give time to mine at least 1 block
	return &LocalDevnet{
		Client:    client,
		closer:    closer,
		cancel:    cancel,
		done:      done,
		numMiners: numMiners,
	}, nil
}

func rpcBuilder(nFull int, storage []int) ([]test.TestNode, []test.TestStorageNode, func(), error) {
	fullApis, storaApis, err := mockSbBuilder(nFull, storage)
	if err != nil {
		return nil, nil, nil, err
	}
	fulls := make([]test.TestNode, nFull)
	storers := make([]test.TestStorageNode, len(storage))

	var closers []func()
	for i, a := range fullApis {
		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filecoin", a)
		go func() {
			if err := http.ListenAndServe(":7777", rpcServer); err != nil {
				if err != http.ErrServerClosed {
					panic(err)
				}
			}
		}()
		var err error
		time.Sleep(time.Second)
		fulls[i].FullNode, _, err = client.NewFullNodeRPC("ws://127.0.0.1:7777", nil)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	for i, a := range storaApis {
		rpcServer := jsonrpc.NewServer()
		rpcServer.Register("Filecoin", a)
		testServ := httptest.NewServer(rpcServer)
		closers = append(closers, func() { testServ.Close() })

		var err error
		storers[i].StorageMiner, _, err = client.NewStorageMinerRPC("ws://"+testServ.Listener.Addr().String(), nil)
		if err != nil {
			return nil, nil, nil, err
		}
		storers[i].MineOne = a.MineOne
	}

	return fulls, storers, func() {
		for _, c := range closers {
			c()
		}
	}, nil
}

const nPreseal = 2

func mockSbBuilder(nFull int, storage []int) ([]test.TestNode, []test.TestStorageNode, error) {
	ctx := context.Background()
	mn := mocknet.New(ctx)

	fulls := make([]test.TestNode, nFull)
	storers := make([]test.TestStorageNode, len(storage))

	var minersPk []crypto.PrivKey
	var minersPid []peer.ID
	for i := 0; i < len(storage); i++ {
		pk, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		pid, err := peer.IDFromPrivateKey(pk)
		if err != nil {
			return nil, nil, err
		}
		minersPk = append(minersPk, pk)
		minersPid = append(minersPid, pid)
	}

	var genbuf bytes.Buffer

	// PRESEAL SECTION, TRY TO REPLACE WITH BETTER IN THE FUTURE
	// TODO: would be great if there was a better way to fake the preseals
	var genms []genesis.Miner
	var genaccs []genesis.Actor
	var maddrs []address.Address
	var presealDirs []string
	var keys []*wallet.Key
	for i := 0; i < len(storage); i++ {
		maddr, err := address.NewIDAddress(genesis2.MinerStart + uint64(i))
		if err != nil {
			return nil, nil, err
		}
		tdir, err := ioutil.TempDir("", "preseal-memgen")
		if err != nil {
			return nil, nil, err
		}
		genm, k, err := sbmock.PreSeal(2048, maddr, nPreseal)
		if err != nil {
			return nil, nil, err
		}
		genm.PeerId = minersPid[0]

		wk, err := wallet.NewKey(*k)
		if err != nil {
			return nil, nil, err
		}

		genaccs = append(genaccs, genesis.Actor{
			Type:    genesis.TAccount,
			Balance: big.NewInt(40000000000),
			Meta:    (&genesis.AccountMeta{Owner: wk.Address}).ActorMeta(),
		})

		keys = append(keys, wk)
		presealDirs = append(presealDirs, tdir)
		maddrs = append(maddrs, maddr)
		genms = append(genms, *genm)
	}
	templ := &genesis.Template{
		Accounts: genaccs,
		Miners:   genms,
	}
	// END PRESEAL SECTION

	for i := 0; i < nFull; i++ {
		var genesis node.Option
		if i == 0 {
			genesis = node.Override(new(modules.Genesis), modtest.MakeGenesisMem(&genbuf, *templ))
		} else {
			genesis = node.Override(new(modules.Genesis), modules.LoadGenesis(genbuf.Bytes()))
		}

		var err error
		// TODO: Don't ignore stop
		_, err = node.New(ctx,
			node.FullAPI(&fulls[i].FullNode),
			node.Online(),
			node.Repo(repo.NewMemory(nil)),
			node.MockHost(mn),
			node.Test(),

			node.Override(new(sectorbuilder.Verifier), sbmock.MockVerifier),

			genesis,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	for i, full := range storage {
		if full != 0 {
			return nil, nil, fmt.Errorf("storage nodes only supported on the first full node")
		}

		f := fulls[full]
		if _, err := f.FullNode.WalletImport(ctx, &keys[i].KeyInfo); err != nil {
			return nil, nil, err
		}
		if err := f.FullNode.WalletSetDefault(ctx, keys[i].Address); err != nil {
			return nil, nil, err
		}

		genMiner := maddrs[i]
		wa := genms[i].Worker

		storers[i] = testStorageNode(ctx, wa, genMiner, minersPk[i], f, mn, node.Options(
			node.Override(new(sealmgr.Manager), func() (sealmgr.Manager, error) {
				return sealmgr.NewSimpleManager(storedcounter.New(datastore.NewMapDatastore(), datastore.NewKey("/potato")), genMiner, sbmock.NewMockSectorBuilder(5, build.SectorSizes[0]))
			}),
			node.Unset(new(*advmgr.Manager)),
		))
	}

	if err := mn.LinkAll(); err != nil {
		return nil, nil, err
	}

	return fulls, storers, nil
}

func testStorageNode(ctx context.Context, waddr address.Address, act address.Address, pk crypto.PrivKey, tnd test.TestNode, mn mocknet.Mocknet, opts node.Option) test.TestStorageNode {
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
	ds, err := lr.Datastore("/metadata")
	if err != nil {
		panic(err)
	}
	err = ds.Put(datastore.NewKey("miner-address"), act.Bytes())
	if err != nil {
		panic(err)
	}
	nic := storedcounter.New(ds, datastore.NewKey("/storage/nextid"))
	for i := 0; i < nPreseal; i++ {
		nic.Next()
	}

	err = lr.Close()
	if err != nil {
		panic(err)
	}
	peerid, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		panic(err)
	}
	enc, err := actors.SerializeParams(&saminer.ChangePeerIDParams{NewID: peerid})
	if err != nil {
		panic(err)
	}
	msg := &types.Message{
		To:       act,
		From:     waddr,
		Method:   builtin.MethodsMiner.ChangePeerID,
		Params:   enc,
		Value:    types.NewInt(0),
		GasPrice: types.NewInt(0),
		GasLimit: 1000000,
	}

	_, err = tnd.MpoolPushMessage(ctx, msg)
	if err != nil {
		panic(err)
	}
	// start node
	var minerapi api.StorageMiner

	mineBlock := make(chan struct{})
	// TODO: use stop
	_, err = node.New(ctx,
		node.StorageMiner(&minerapi),
		node.Online(),
		node.Repo(r),
		node.Test(),

		node.MockHost(mn),

		node.Override(new(api.FullNode), tnd),
		node.Override(new(*miner.Miner), miner.NewTestMiner(mineBlock, act)),

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
	mineOne := func(ctx context.Context) error {
		select {
		case mineBlock <- struct{}{}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return test.TestStorageNode{StorageMiner: minerapi, MineOne: mineOne}
}
