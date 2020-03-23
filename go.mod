module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-sectorbuilder v0.0.2-0.20200203173614-42d67726bb62
	github.com/filecoin-project/lotus v0.2.11-0.20200320211838-628a598ca064
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-ipfs-blockstore v0.1.4 // indirect
	github.com/ipfs/go-ipld-cbor v0.0.5-0.20200204214505-252690b78669 // indirect
	github.com/ipfs/go-log v1.0.2 // indirect
	github.com/ipfs/go-log/v2 v2.0.2
	github.com/ipfs/go-merkledag v0.3.1 // indirect
	github.com/libp2p/go-libp2p v0.5.2
	github.com/libp2p/go-libp2p-core v0.4.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
	github.com/textileio/lotus-client v0.0.0-20200323231235-5c42584d6d6d
	github.com/whyrusleeping/cbor-gen v0.0.0-20200223203819-95cdfde1438f // indirect
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
