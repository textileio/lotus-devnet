module github.com/textileio/lotus-devnet

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.0
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/filecoin-project/go-address v0.0.3
	github.com/filecoin-project/go-fil-markets v0.6.2
	github.com/filecoin-project/go-jsonrpc v0.1.2-0.20200822201400-474f4fdccc52
	github.com/filecoin-project/go-state-types v0.0.0-20200911004822-964d6c679cfc
	github.com/filecoin-project/go-storedcounter v0.0.0-20200421200003-1c99c62e8a5b
	github.com/filecoin-project/lotus v0.7.2
	github.com/filecoin-project/sector-storage v0.0.0-20200730050024-3ee28c3b6d9a // indirect
	github.com/filecoin-project/specs-actors v0.9.10
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log/v2 v2.1.2-0.20200626104915-0016c0b4b3e4
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.6.1
	go.uber.org/dig v1.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect

)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi

replace github.com/supranational/blst => github.com/supranational/blst v0.1.2-alpha.1

//replace github.com/filecoin-project/lotus => github.com/jsign/lotus v0.4.1-0.20200925143422-2f0854cad509
