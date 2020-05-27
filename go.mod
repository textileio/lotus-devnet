module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-fil-markets v0.2.7
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200520183639-7c6ee2e066b4
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.3.1-0.20200527125918-af1d54501fa6
	github.com/filecoin-project/sector-storage v0.0.0-20200522011946-a59ca7536a95
	github.com/filecoin-project/specs-actors v0.5.4-0.20200521014528-0df536f7e461
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.0.8
	github.com/libp2p/go-libp2p v0.9.2
	github.com/libp2p/go-libp2p-core v0.5.6
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
