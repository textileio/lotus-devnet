module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/filecoin-project/go-address v0.0.5-0.20201103152444-f2023ef3f5bb
	github.com/filecoin-project/go-fil-markets v1.0.4
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20201008195726-68c6a2704e49
	github.com/filecoin-project/go-state-types v0.0.0-20201013222834-41ea465f274f
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v1.1.3
	github.com/filecoin-project/specs-actors v0.9.12
	github.com/filecoin-project/specs-actors/v2 v2.2.0
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4
	github.com/libp2p/go-libp2p v0.12.0
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect

)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

replace github.com/supranational/blst => ./extern/blst

replace github.com/filecoin-project/lotus => github.com/jsign/lotus v0.4.1-0.20201113110440-4640b443453c
