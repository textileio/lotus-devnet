package test

import (
	"context"
	"fmt"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/test"
	sealing "github.com/filecoin-project/lotus/extern/storage-sealing"
)

func StartSealingWaiting(ctx context.Context, miner test.TestStorageNode) {
	snums, err := miner.SectorsList(ctx)
	if err != nil {
		panic(err)
	}

	for _, snum := range snums {
		si, err := miner.SectorsStatus(ctx, snum, false)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Sector state: %s", si.State)
		if si.State == api.SectorState(sealing.WaitDeals) {
			if err := miner.SectorStartSealing(ctx, snum); err != nil {
				panic(err)
			}
		}
	}
}
