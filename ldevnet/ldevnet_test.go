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
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/lib/jsonrpc"
	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"
	"github.com/textileio/lotus-client/api"
	"github.com/textileio/lotus-client/api/apistruct"
)

func TestMain(m *testing.M) {
	logging.SetAllLoggers(logging.LevelDebug)
	os.Exit(m.Run())
}

func TestStore(t *testing.T) {
	_, err := New(1, time.Millisecond*50)
	require.Nil(t, err)

	var client apistruct.FullNodeStruct
	_, err = jsonrpc.NewMergeClient("ws://127.0.0.1:7777/rpc/v0", "Filecoin",
		[]interface{}{
			&client.Internal,
			&client.CommonStruct.Internal,
		}, nil)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	ctx := context.Background()

	ts, err := client.ChainHead(ctx)
	require.Nil(t, err)
	require.Greater(t, int64(ts.Height()), int64(5))

	miners, err := client.StateListMiners(ctx, types.EmptyTSK)
	require.Nil(t, err)
	require.Greater(t, len(miners), 0)

	waddr, err := client.WalletDefaultAddress(ctx)
	require.Nil(t, err)
	require.NotEmpty(t, waddr)

	tmpf, err := ioutil.TempFile("", "")
	require.Nil(t, err)

	data := make([]byte, 600)
	rand.New(rand.NewSource(5)).Read(data)
	err = ioutil.WriteFile(tmpf.Name(), data, 0644)
	require.Nil(t, err)
	fcid, err := client.ClientImport(ctx, api.FileRef{Path: tmpf.Name()})
	require.Nil(t, err)
	require.True(t, fcid.Defined())

	sdp := &api.StartDealParams{
		Data:              &storagemarket.DataRef{Root: fcid},
		Wallet:            waddr,
		EpochPrice:        types.NewInt(1000000),
		MinBlocksDuration: 100,
		Miner:             miners[0],
	}
	deal, err := client.ClientStartDeal(ctx, sdp)
	require.Nil(t, err)

	time.Sleep(time.Second)
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
		fmt.Println(di.State)
		time.Sleep(time.Second)
	}

	offers, err := client.ClientFindData(ctx, fcid)
	require.Nil(t, err)
	require.Greater(t, len(offers), 0)

	rpath, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(rpath)

	ref := api.FileRef{
		Path:  filepath.Join(rpath, "ret"),
		IsCAR: false,
	}
	err = client.ClientRetrieve(ctx, offers[0].Order(waddr), ref)
	require.Nil(t, err)

	rdata, err := ioutil.ReadFile(filepath.Join(rpath, "ret"))
	require.Nil(t, err)
	require.True(t, bytes.Equal(data, rdata))
}
