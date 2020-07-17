module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/go-fil-markets v0.4.0
	github.com/filecoin-project/go-jsonrpc v0.1.1-0.20200602181149-522144ab4e24
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.4.2-0.20200710025032-337c57093404
	github.com/filecoin-project/sector-storage v0.0.0-20200709184611-f0dae546b517
	github.com/filecoin-project/specs-actors v0.7.1
	github.com/filecoin-project/storage-fsm v0.0.0-20200707194229-bc5e298e2b4c
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

// Includes a fix of go-fil-markets choosing wrong wallet address
replace github.com/filecoin-project/go-fil-markets => github.com/jsign/go-fil-markets v0.4.1-0.20200716182644-ae4ece2facad

// Includes a fix of lotus about using IPFS as the blockstore
replace github.com/filecoin-project/lotus => github.com/jsign/lotus v0.4.1-0.20200717230639-a1439f33f32c
