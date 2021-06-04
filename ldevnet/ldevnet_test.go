package ldevnet

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc"
	lapi "github.com/filecoin-project/lotus/api"
	api "github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	//logging.SetAllLoggers(logging.LevelDebug)
	//logging.SetLogLevel("miner", "ERROR")
	//logging.SetLogLevel("chainstore", "ERROR")
	//logging.SetLogLevel("chain", "ERROR")
	//logging.SetLogLevel("sub", "ERROR")
	//logging.SetLogLevel("storageminer", "ERROR")
	_ = logging.ErrNoSuchLogger
	os.Exit(m.Run())
}

func TestStore(t *testing.T) {
	numMiners := []int{1}

	for _, nm := range numMiners {
		for i := 0; i < 1; i++ {
			t.Run(fmt.Sprintf("%d miners, deal with miner %d", nm, i), dealSpecificMiner(t, nm, 0))
		}
	}
}

func dealSpecificMiner(t *testing.T, numMiners int, concreteMiner int) func(*testing.T) {
	return func(t *testing.T) {
		_, err := New(numMiners, time.Millisecond*100, true, "", false)
		require.Nil(t, err)

		var client api.FullNodeStruct
		cc, err := jsonrpc.NewMergeClient(context.Background(), "ws://127.0.0.1:7777/rpc/v0", "Filecoin",
			[]interface{}{
				&client.Internal,
				&client.CommonStruct.Internal,
			}, nil)
		if err != nil {
			panic(err)
		}
		defer cc()
		time.Sleep(time.Second)
		ctx := context.Background()

		ts, err := client.ChainHead(ctx)
		require.Nil(t, err)
		require.Greater(t, int64(ts.Height()), int64(2))

		miners, err := client.StateListMiners(ctx, types.EmptyTSK)
		require.Nil(t, err)
		require.Greater(t, len(miners), 0)

		waddr, err := client.WalletDefaultAddress(ctx)
		require.Nil(t, err)
		require.NotEmpty(t, waddr)

		tmpf, err := ioutil.TempFile("", "")
		require.Nil(t, err)

		r := rand.New(rand.NewSource(22))
		for i := 0; i < 1; i++ {
			data := make([]byte, 100*1024*1024)
			r.Read(data)
			err = ioutil.WriteFile(tmpf.Name(), data, 0644)
			require.Nil(t, err)
			fcid, err := client.ClientImport(ctx, lapi.FileRef{Path: tmpf.Name()})
			require.Nil(t, err)
			require.True(t, fcid.Root.Defined())

			sdp := &lapi.StartDealParams{
				Data: &storagemarket.DataRef{
					TransferType: storagemarket.TTGraphsync,
					Root:         fcid.Root,
				},
				Wallet:            waddr,
				EpochPrice:        types.NewInt(100000000),
				MinBlocksDuration: uint64(build.MinDealDuration),
				Miner:             miners[concreteMiner],
			}
			deal, err := client.ClientStartDeal(ctx, sdp)
			require.Nil(t, err)

			time.Sleep(time.Second / 2)

		loop:
			for {
				di, err := client.ClientGetDealInfo(ctx, *deal)
				require.Nil(t, err)

				switch di.State {
				case storagemarket.StorageDealProposalRejected:
					t.Fatal("deal rejected")
				case storagemarket.StorageDealFailing:
					t.Fatal("deal failed")
				case storagemarket.StorageDealError:
					t.Fatal("deal errored")
				case storagemarket.StorageDealActive:
					fmt.Println("COMPLETE", di)
					break loop
				}
				fmt.Println(storagemarket.DealStates[di.State])
				time.Sleep(time.Second)
			}
			offers, err := client.ClientFindData(ctx, fcid.Root, nil)
			require.Nil(t, err)
			require.Greater(t, len(offers), 0)

			rpath, err := ioutil.TempDir("", "")
			require.Nil(t, err)
			defer os.RemoveAll(rpath)

			// Retrieve to file.
			ref := &lapi.FileRef{
				Path:  filepath.Join(rpath, "ret"),
				IsCAR: false,
			}
			err = client.ClientRetrieve(ctx, offers[0].Order(waddr), ref)
			require.Nil(t, err)

			rdata, err := ioutil.ReadFile(filepath.Join(rpath, "ret"))
			require.Nil(t, err)
			require.True(t, bytes.Equal(data, rdata))
		}
	}
}
