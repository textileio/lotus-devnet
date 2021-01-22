package test

import (
	"context"
	"fmt"
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/test"
	"github.com/filecoin-project/lotus/miner"
)

func MineUntilBlock(ctx context.Context, fn test.TestNode, sn test.TestStorageNode, cb func(abi.ChainEpoch)) {
	for i := 0; i < 1000; i++ {
		var success bool
		var err error
		var epoch abi.ChainEpoch
		wait := make(chan struct{})
		mineErr := sn.MineOne(ctx, miner.MineReq{
			Done: func(win bool, ep abi.ChainEpoch, e error) {
				success = win
				err = e
				epoch = ep
				wait <- struct{}{}
			},
		})
		if mineErr != nil {
			panic(mineErr)
		}
		<-wait
		if err != nil {
			panic(err)
		}
		if success {
			// Wait until it shows up on the given full nodes ChainHead
			nloops := 50
			for i := 0; i < nloops; i++ {
				ts, err := fn.ChainHead(ctx)
				if err != nil {
					panic(err)
				}
				if ts.Height() == epoch {
					break
				}
				if i == nloops-1 {
					panic(err)
				}
				time.Sleep(time.Millisecond * 10)
			}

			if cb != nil {
				cb(epoch)
			}
			return
		}
		fmt.Printf("did not mine block, trying again %v", i)
	}
	panic("failed to mine 1000 times in a row...")
}
