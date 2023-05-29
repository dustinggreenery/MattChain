package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/dustinggreenery/MattChain/utils"
)

type MattChain struct {
	LastHash []byte
	Database *badger.DB
};

func InitBlockChain(address string) *MattChain {	
	if utils.DBexists() {
		fmt.Println("Blockchain Exists Already");
		runtime.Goexit();
	}

	var lastHash []byte;

	db, err := badger.Open(badger.DefaultOptions(utils.DbPath));
	utils.ErrorHandling(err);

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, utils.GenesisData);
		genesis := Genesis(cbtx);
		fmt.Println("Genesis Created");

		err = txn.Set(genesis.Hash, genesis.SerializeBlock());
		utils.ErrorHandling(err);
		
		err = txn.Set([]byte("lh"), genesis.Hash);
		lastHash = genesis.Hash;

		return err;
	});
	utils.ErrorHandling(err);

	mattchain := MattChain{lastHash, db};
	return &mattchain;
}

func ContinueBlockChain(address string) *MattChain {
	if !utils.DBexists() {
		fmt.Println("No Existing Blockchain");
		runtime.Goexit();
	}

	var lastHash []byte;

	db, err := badger.Open(badger.DefaultOptions(utils.DbPath));
	utils.ErrorHandling(err);

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"));
		utils.ErrorHandling(err);

		err = item.Value(func(val []byte) error {
			lastHash = val;
			return nil;
		});

		return err;
	})
	utils.ErrorHandling(err);

	chain := MattChain{lastHash, db};
	return &chain;
}

func (chain *MattChain) AddBlock(transactions []*MattTransaction) *MattBlock {
	var lastHash []byte;

	for _, tx := range transactions {
		if !chain.VerifyTransaction(tx) {
			log.Panic("Invalid Transaction");
		}
	}

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"));
		utils.ErrorHandling(err);

		err = item.Value(func(val []byte) error {
			lastHash = val;
			return nil;
		})

		return err;
	})
	utils.ErrorHandling(err);

	newBlock := CreateBlock(transactions, lastHash);

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.SerializeBlock());
		utils.ErrorHandling(err);

		err = txn.Set([]byte("lh"), newBlock.Hash);
		chain.LastHash = newBlock.Hash;

		return err;
	});

	utils.ErrorHandling(err);

	return newBlock;
}

func (chain *MattChain) Iterator() *MattChainIterator {
	iter := &MattChainIterator{chain.LastHash, chain.Database};
	return iter;
}

func (chain *MattChain) FindTransaction(ID []byte) (MattTransaction, error) {
	iter := chain.Iterator();

	for {
		block := iter.Next();

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil;
			}
		}

		if len(block.PrevHash) == 0 {
			break;
		}
	}

	return MattTransaction{}, errors.New("Transaction doesn't exist");
}

func (chain *MattChain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs);
	spentTXOs := make(map[string][]int);
	iter := chain.Iterator();

	for {
		block := iter.Next();

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID);

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs;
						}
					}
				}

				outs := UTXO[txID];
				outs.Outputs = append(outs.Outputs, out);
				UTXO[txID] = outs;
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID);
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out);
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break;
		}
	}
	return UTXO;
}

func (chain *MattChain) SignTransaction(tx *MattTransaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]MattTransaction);
	
	for _, in := range tx.Inputs {
		prevTx, err := chain.FindTransaction(in.ID);
		utils.ErrorHandling(err);

		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx;
	}

	tx.Sign(privKey, prevTXs);
}

func (chain *MattChain) VerifyTransaction(tx *MattTransaction) bool {
	if tx.IsCoinbase() {
		return true;
	}

	prevTXs := make(map[string]MattTransaction);
	
	for _, in := range tx.Inputs {
		prevTx, err := chain.FindTransaction(in.ID);
		utils.ErrorHandling(err);
		
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx;
	}

	return tx.Verify(prevTXs);
}
