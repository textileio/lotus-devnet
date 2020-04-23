module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.0.0-20200415011556-4378bd41b91f
	github.com/filecoin-project/lotus v0.2.11-0.20200423065310-a99875feaca9
	github.com/filecoin-project/sector-storage v0.0.0-20200417225459-e75536581a08
	github.com/filecoin-project/specs-actors v0.0.0-20200421235624-312ac81e2aa4
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.3
	github.com/libp2p/go-libp2p v0.6.1
	github.com/libp2p/go-libp2p-core v0.5.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/textileio/lotus-client v0.0.0-20200423150613-e9939eba3b2b
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
