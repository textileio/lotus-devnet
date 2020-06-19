module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-fil-markets v0.3.0
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200602181149-522144ab4e24
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.4.0
	github.com/filecoin-project/sector-storage v0.0.0-20200618073200-d9de9b7cb4b4
	github.com/filecoin-project/specs-actors v0.6.2-0.20200617175406-de392ca14121
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.1.2-0.20200609205458-f8d20c392cb7
	github.com/libp2p/go-libp2p v0.9.4
	github.com/libp2p/go-libp2p-core v0.5.7
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.5.1
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
