module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-fil-markets v0.5.2
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200602181149-522144ab4e24
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.4.3-0.20200729021736-477dd536b5b5
	github.com/filecoin-project/sector-storage v0.0.0-20200727112136-9377cb376d25
	github.com/filecoin-project/specs-actors v0.8.1-0.20200728182452-1476088f645b
	github.com/filecoin-project/storage-fsm v0.0.0-20200728185042-33f96f051f20
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4
	github.com/libp2p/go-libp2p v0.10.0
	github.com/libp2p/go-libp2p-core v0.6.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543

)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

// Includes a fix for multiple pieces in a sector in a mocked sectorbuilder.
// jsign/9377cbpatched = sector-storage@9777cbpatched plus a patch commit for big sectors.
replace github.com/filecoin-project/sector-storage => github.com/jsign/sector-storage v0.0.0-20200730134255-570a28cf3edf
