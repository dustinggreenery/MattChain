package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dustinggreenery/MattChain/blockchain"
	"github.com/dustinggreenery/MattChain/utils"
	"github.com/dustinggreenery/MattChain/wallet"
)

func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets();
	address := wallets.AddWallet();
	wallets.SaveFile();

	fmt.Printf("New address is: %s\n", address);
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets();
	addresses := wallets.GetAllAddresses();

	for _, address := range addresses {
		fmt.Println(address);
	}
}

func (cli *CommandLine) getBalance(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not valid");
	}

	chain := blockchain.ContinueBlockChain(address);
	UTXOSet := blockchain.UTXOSet{chain};
	defer chain.Database.Close();

	balance := 0;
	pubKeyHash := utils.Base58Decode([]byte(address));
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4];
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash);

	for _, out := range UTXOs {
		balance += out.Value;
	}

	fmt.Printf("Balance of %s: %d\n", address, balance);
}


func (cli *CommandLine) createBlockChain(address string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("Address isn't Valid");
	}

	chain := blockchain.InitBlockChain(address);
	defer chain.Database.Close();

	UTXOSet := blockchain.UTXOSet{chain};
	UTXOSet.Reindex();

	fmt.Println("Finished!");
}

func (cli *CommandLine) send(from, to string, amount int) {
	if !wallet.ValidateAddress(to) {
		log.Panic("Sending Address isn't Valid");
	}
	
	if !wallet.ValidateAddress(from) {
		log.Panic("Receiving Address isn't Valid");
	}

	chain := blockchain.ContinueBlockChain(from);
	UTXOSet := blockchain.UTXOSet{chain};
	defer chain.Database.Close();

	tx := blockchain.NewTransaction(from, to, amount, &UTXOSet);
	cbTx := blockchain.CoinbaseTx(from, "");
	block := chain.AddBlock([]*blockchain.MattTransaction{cbTx, tx});

	UTXOSet.Update(block);
	
	fmt.Println("Success!");
}

func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockChain("");
	defer chain.Database.Close();
	iter := chain.Iterator();

	for {
		block := iter.Next();

		fmt.Printf("Hash: %x\n", block.Hash);
		fmt.Printf("  Previous Hash: %x\n", block.PrevHash);

		pow := blockchain.NewProof(block);
		fmt.Printf("  PoW: %s\n", strconv.FormatBool(pow.Validate()));

		for _, tx := range block.Transactions {
			fmt.Println(tx);
		}

		fmt.Println();

		if len(block.PrevHash) == 0 {
			break;
		}
	}
}

func (cli *CommandLine) reindexUTXO() {
	chain := blockchain.ContinueBlockChain("");
	defer chain.Database.Close();
	UTXOSet := blockchain.UTXOSet{chain};
	UTXOSet.Reindex();

	count := UTXOSet.CountTransactions();
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count);
}