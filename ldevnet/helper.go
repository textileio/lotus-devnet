package ldevnet

import (
	"context"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
)

func init() {
	build.InsecurePoStValidation = true
}

func CreateLocalDevnet(t *testing.T, numMiners int, bigSectors bool) (*LocalDevnet, address.Address, []address.Address, func()) {
	dnet, err := New(numMiners, DefaultDuration, bigSectors)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	addr, err := dnet.Client.WalletDefaultAddress(ctx)
	if err != nil {
		t.Fatal(err)
	}

	miners, err := dnet.Client.StateListMiners(ctx, types.EmptyTSK)
	if err != nil {
		t.Fatal(err)
	}

	return dnet, addr, miners, func() { dnet.Close() }
}
