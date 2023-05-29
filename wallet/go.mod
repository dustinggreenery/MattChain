module github.com/dustinggreenery/MattChain/wallet

go 1.18

require (
	github.com/dustinggreenery/MattChain/utils v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.8.0
)

require github.com/mr-tron/base58 v1.2.0 // indirect

replace github.com/dustinggreenery/MattChain/utils => ../utils
