package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dustinggreenery/MattChain/utils"
)

const walletFile = "../tmp/wallets.data"

type MattWallets struct {
	Wallets map[string]*MattWallet
}

func CreateWallets() (*MattWallets, error) {
	wallets := MattWallets{};
	wallets.Wallets = make(map[string]*MattWallet);

	err := wallets.LoadFile();

	return &wallets, err;
}

func (ws *MattWallets) AddWallet() string {
	wallet := MakeWallet();
	address := fmt.Sprintf("%s", wallet.Address());

	ws.Wallets[address] = wallet;

	return address;
}

func (ws MattWallets) GetWallet(address string) MattWallet {
	return *ws.Wallets[address];
}

func (ws *MattWallets) GetAllAddresses() []string {
	var addresses []string;

	for address := range ws.Wallets {
		addresses = append(addresses, address);
	}

	return addresses;
}

func (ws *MattWallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err;
	}

	var wallets MattWallets;

	fileContent, err := ioutil.ReadFile(walletFile);
	if err != nil {
		return err;
	}

	gob.Register(elliptic.P256());

	decoder := gob.NewDecoder(bytes.NewReader(fileContent));
	err = decoder.Decode(&wallets);
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets;

	return nil
}

func (ws *MattWallets) SaveFile() {
	var content bytes.Buffer;

	gob.Register(elliptic.P256());

	encoder := gob.NewEncoder(&content);
	err := encoder.Encode(ws);
	utils.ErrorHandling(err);

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644);
	utils.ErrorHandling(err);
}