package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/dustinggreenery/MattChain/utils"
)

type CommandLine struct {}

func (cli *CommandLine) printUsage() {
	fmt.Println("Commands:");
	
	fmt.Println(" createwallet - Creates a new wallet with an address.");
	fmt.Println(" listaddresses - Lists all addresses of the wallets created.");
	fmt.Println(" getbalance - Gets the balance of a wallet.");
	fmt.Println("	-address Address of wallet");

	fmt.Println();
	fmt.Println(" createblockchain - Creates a blockchain if not created already.");
	fmt.Println("	-address Address that mines coinbase transaction who also gets the genesis reward.");
	fmt.Println(" send - Send tokens from one address to another.");
	fmt.Println("	-from Address sending funds");
	fmt.Println("	-to Address receiving funds");
	fmt.Println("	-amount Amount of tokens to be sent");
	fmt.Println(" printchain - Prints the blocks and their transactions on the chain.");
	fmt.Println(" reindexutxo - Reindexes the unspent transaction outputs.");
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage();
		runtime.Goexit();
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs();


	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError);
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError);
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError);

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError);
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError);
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError);
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError);

	getBalanceAddress := getBalanceCmd.String("address", "", "Address to get balance of");
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Address that mines coinbase transaction who also gets the genesis reward.");
	sendFrom := sendCmd.String("from", "", "Address sending funds");
	sendTo := sendCmd.String("to", "", "Address receiving funds");
	sendAmount := sendCmd.Int("amount", 0, "Amount of tokens to be sent");
	
	
	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "send":
		err := sendCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:]);
		utils.ErrorHandling(err);
	default:
		cli.printUsage();
		runtime.Goexit();
	}


	if createWalletCmd.Parsed() {
		cli.createWallet();
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses();
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage();
			runtime.Goexit();
		}

		cli.getBalance(*getBalanceAddress);
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage();
			runtime.Goexit();
		}

		cli.createBlockChain(*createBlockchainAddress);
	}
	
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage();
			runtime.Goexit();
		}

		cli.send(*sendFrom, *sendTo, *sendAmount);
	}

	if printChainCmd.Parsed() {
		cli.printChain();
	}
	
	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO();
	}
}