package ldevnet

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"
	"github.com/textileio/lotus-api/api"
	"github.com/textileio/lotus-api/api/apistruct"
	"github.com/textileio/lotus-api/chain/types"
	"github.com/textileio/lotus-api/lib/jsonrpc"
)

func TestMain(m *testing.M) {
	logging.SetAllLoggers(logging.LevelError)
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
	require.Greater(t, ts.Height(), uint64(5))

	miners, err := client.StateListMiners(ctx, types.EmptyTSK)
	require.Nil(t, err)
	require.Greater(t, len(miners), 0)

	waddr, err := client.WalletDefaultAddress(ctx)
	require.Nil(t, err)
	require.NotEmpty(t, waddr)

	tmpf, err := ioutil.TempFile("", "")
	require.Nil(t, err)

	data := make([]byte, 1000)
	rand.New(rand.NewSource(22)).Read(data)
	err = ioutil.WriteFile(tmpf.Name(), data, 0644)
	require.Nil(t, err)
	fcid, err := client.ClientImport(ctx, tmpf.Name())
	require.Nil(t, err)
	require.True(t, fcid.Defined())

	deal, err := client.ClientStartDeal(ctx, fcid, waddr, miners[0], types.NewInt(40000000), 100)
	require.Nil(t, err)

	time.Sleep(time.Second)
loop:
	for {
		di, err := client.ClientGetDealInfo(ctx, *deal)
		require.Nil(t, err)

		switch di.State {
		case api.DealRejected:
			t.Fatal("deal rejected")
		case api.DealFailed:
			t.Fatal("deal failed")
		case api.DealError:
			t.Fatal("deal errored")
		case api.DealComplete:
			break loop
		}
		time.Sleep(time.Second)
	}

	offers, err := client.ClientFindData(ctx, fcid)
	require.Nil(t, err)
	require.Greater(t, len(offers), 0)

	rpath, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(rpath)

	err = client.ClientRetrieve(ctx, offers[0].Order(waddr), filepath.Join(rpath, "ret"))
	require.Nil(t, err)

	rdata, err := ioutil.ReadFile(filepath.Join(rpath, "ret"))
	require.Nil(t, err)
	require.Equal(t, data, rdata)
}
