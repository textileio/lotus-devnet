package ldevnet

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/test"
	"github.com/filecoin-project/lotus/chain/actors/policy"

	api_test_internal "github.com/textileio/lotus-devnet/ldevnet/api/test"
	node_internal "github.com/textileio/lotus-devnet/ldevnet/node"
	test_internal "github.com/textileio/lotus-devnet/ldevnet/node/test"
)

// Related files:
// - node/node_test.go
// - node/test/builder.go
// - api/test/test.go
// - api/test/deals.go

const (
	DefaultDurationMs = 100
	DefaultDuration   = time.Millisecond * DefaultDurationMs
)

func init() {
	node_internal.Init()

	os.Setenv("TRUST_PARAMS", "1")
	os.Setenv("LOTUS_DISABLE_WATCHDOG", "1")
}

type LocalDevnet struct {
	Client *api.FullNodeStruct

	numMiners int
	cancel    context.CancelFunc
	done      chan struct{}
}

func (ld *LocalDevnet) Close() {
	ld.cancel()

	for i := 0; i < ld.numMiners; i++ {
		<-ld.done
	}
	close(ld.done)
}

func New(numMiners int, blockDur time.Duration, bigSector bool, ipfsAddr string, onlineMode bool) (*LocalDevnet, error) {
	if bigSector {
		policy.SetSupportedProofTypes(abi.RegisteredSealProof_StackedDrg512MiBV1)
	}

	// Create full-node and miners
	miners := make([]test.StorageMiner, numMiners)
	for i := 0; i < numMiners; i++ {
		miners[i] = test.OneMiner[0]
	}
	n, sn, err := rpcBuilder(miners, bigSector, ipfsAddr, onlineMode)
	if err != nil {
		return nil, err
	}

	// Connect everyone.
	client := n.FullNode.(*api.FullNodeStruct)
	ctx, cancel := context.WithCancel(context.Background())
	addrinfo, err := client.NetAddrsListen(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	for i := range miners {
		if err := sn[i].NetConnect(ctx, addrinfo); err != nil {
			cancel()
			return nil, err
		}
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
		time.Sleep(time.Second)
	}

	// Mine.
	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()
		for {
			time.Sleep(blockDur)
			if ctx.Err() != nil {
				panic(ctx.Err())
			}
			if err := sn[0].MineOne(context.Background(), test.MineNext); err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		for {
			for i := 0; i < numMiners; i++ {
				api_test_internal.StartSealingWaiting(ctx, sn[i])
			}
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(blockDur * 2) // Give time to mine at least 1 block
	return &LocalDevnet{
		Client:    client,
		cancel:    cancel,
		done:      done,
		numMiners: numMiners,
	}, nil
}

func rpcBuilder(storage []test.StorageMiner, bigSector bool, ipfsAddr string, onlineMode bool) (test.TestNode, []test.TestStorageNode, error) {
	n, sn, err := test_internal.MockSbBuilder(test.OneFull, storage, bigSector, ipfsAddr, onlineMode)
	if err != nil {
		return test.TestNode{}, nil, err
	}

	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("Filecoin", n[0])
	go func() {
		if err := http.ListenAndServe(":7777", rpcServer); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()
	time.Sleep(time.Second)

	n[0].FullNode, _, err = client.NewFullNodeRPCV1(context.Background(), "ws://127.0.0.1:7777", nil)
	if err != nil {
		return test.TestNode{}, nil, err
	}

	return n[0], sn, nil
}
