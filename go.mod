module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200218010043-eb9bb40ed5be
	github.com/filecoin-project/go-fil-markets v0.1.0
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.2.11-0.20200429011108-559b5e788523
	github.com/filecoin-project/sector-storage v0.0.0-20200425102315-c19a25449861
	github.com/filecoin-project/specs-actors v0.2.1-0.20200428232403-f0282340f59a
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.5
	github.com/libp2p/go-libp2p v0.8.1
	github.com/libp2p/go-libp2p-core v0.5.1
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
