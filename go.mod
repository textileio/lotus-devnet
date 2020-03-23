module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.0.0-20200318012938-6403a5bda668
	github.com/filecoin-project/go-sectorbuilder v0.0.2-0.20200314022627-38af9db49ba2
	github.com/filecoin-project/lotus v0.2.11-0.20200321025511-d1d2e753d946
	github.com/filecoin-project/specs-actors v0.0.0-20200312030511-3f5510bf6130
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.2
	github.com/textileio/lotus-api v0.0.0-20200323185520-61a0ab7984f0
	github.com/libp2p/go-libp2p v0.5.2
	github.com/libp2p/go-libp2p-core v0.4.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
