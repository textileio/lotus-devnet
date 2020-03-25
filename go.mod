module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.0.0-20200318012938-6403a5bda668
	github.com/filecoin-project/go-sectorbuilder v0.0.2-0.20200317221918-42574fc2aab9
	github.com/filecoin-project/lotus v0.2.11-0.20200325000223-ca3d2bf46f69
	github.com/filecoin-project/specs-actors v0.0.0-20200321055844-54fa2e8da1c2
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.2
	github.com/libp2p/go-libp2p v0.6.0
	github.com/libp2p/go-libp2p-core v0.5.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
