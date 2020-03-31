module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.0.0-20200318012938-6403a5bda668
	github.com/filecoin-project/go-sectorbuilder v0.0.2-0.20200317221918-42574fc2aab9
	github.com/filecoin-project/lotus v0.2.11-0.20200331202907-e94d136e3558
	github.com/filecoin-project/specs-actors v0.0.0-20200324235424-aef9b20a9fb1
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.3
	github.com/libp2p/go-libp2p v0.6.0
	github.com/libp2p/go-libp2p-core v0.5.0
	github.com/prometheus/common v0.4.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
	github.com/textileio/lotus-client v0.0.0-20200331223428-5842b040af61
	go4.org v0.0.0-20190313082347-94abd6928b1d // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
