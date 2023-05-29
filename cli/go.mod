module github.com/dustinggreenery/MattChain/cli

go 1.18

replace github.com/dustinggreenery/MattChain/wallet => ../wallet

replace github.com/dustinggreenery/MattChain/blockchain => ../blockchain

require (
	github.com/dustinggreenery/MattChain/blockchain v0.0.0-00010101000000-000000000000
	github.com/dustinggreenery/MattChain/utils v0.0.0-00010101000000-000000000000
	github.com/dustinggreenery/MattChain/wallet v0.0.0-00010101000000-000000000000
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dustinggreenery/MattChain/merkle v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/glog v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/dustinggreenery/MattChain/utils => ../utils

replace github.com/dustinggreenery/MattChain/merkle => ../merkle
