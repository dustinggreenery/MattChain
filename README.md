# Mattchain

![Contributors](https://img.shields.io/github/contributors/dustinggreenery/MattChain)
![Forks](https://img.shields.io/github/forks/dustinggreenery/MattChain)
![Stars](https://img.shields.io/github/stars/dustinggreenery/MattChain)
![Licence](https://img.shields.io/github/license/dustinggreenery/MattChain)
![Issues](https://img.shields.io/github/issues/dustinggreenery/MattChain)

### Heavily inspired by Tensor Programming's Golang Blockchain series!

## Description

Mattchain is a blockchain I created using Golang!

### Aspects and Concepts

- Proof of Work
- Persistence
- Transactions
- Wallets
- Signatures
- Merkle Trees

## Usage

### To Get Started

- Fork the repository
- Clone the forked repository
- Go to the file location on a command prompt (.\MattChain\main)
- Make sure you have Golang downloaded
- Use the below command to see your options!

```bash
go run main.go
```

## Commands:

```bash
go run main.go createwallet
```

Creates a new wallet and gives you the public address of the wallet. With this wallet you can create a blockchain and send tokens.

<br />

```bash
go run main.go listaddresses
```

List the address of all the wallets created on your system so far.

<br />

```bash
go run main.go getbalance -address ADDRESS
```

Lists the amount of tokens a wallet has, with a given address.

<br />

```bash
go run main.go createblockchain -address ADDRESS
```

Creates a blockchain using the given address if it wasn't created already. The coinbase transaction prize is also sent to this address' wallet.

<br />

```bash
go run main.go send -from ADDRESS -to ADDRESS -amount AMOUNT
```

Sends tokens from one of your wallets to another.

<br />

```bash
go run main.go printchain
```

Prints the whole chain create so far, with its blocks and stored transactions.

<br />

```bash
go run main.go reindexutxo
```

Reindexes all the unspent transaction outputs, setting the amount of tokens in the wallet to the most current blockchain. The program normally does this for you though.
