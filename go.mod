module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.0.0-20200415011556-4378bd41b91f
	github.com/filecoin-project/lotus v0.2.11-0.20200425000157-1bfa2311d693
	github.com/filecoin-project/specs-actors v1.0.1-0.20200424174946-11410d0bbcaf
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.5
	github.com/libp2p/go-libp2p v0.8.1
	github.com/libp2p/go-libp2p-core v0.5.1
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
