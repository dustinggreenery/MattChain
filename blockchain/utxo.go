package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dustinggreenery/MattChain/utils"
)

var utxoPrefix = []byte("utxo-")

type UTXOSet struct {
	Mattchain *MattChain
}

func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := u.Mattchain.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
		return nil;
	}

	collectSize := 100000;

	u.Mattchain.Database.View(func (txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions;
		opts.PrefetchValues = false;

		it := txn.NewIterator(opts);
		defer it.Close();

		keysForDelete := make([][]byte, 0 , collectSize);
		keysCollected := 0;

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil);
			keysForDelete = append(keysForDelete, key);
			keysCollected++;

			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err);
				}

				keysForDelete = make([][]byte, 0, collectSize);
				keysCollected = 0;
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err);
			}
		}

		return nil;
	})
}

func (u *UTXOSet) FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput;

	db := u.Mattchain.Database;

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions;

		it := txn.NewIterator(opts);
		defer it.Close();

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item();
			
			var v []byte;

			err := item.Value(func(val []byte) error {
				v = val;
				return nil;
			});
			utils.ErrorHandling(err);

			outs := DeserializeOutputs(v);
			
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out);
				}
			}
		}

		return nil;
	});
	utils.ErrorHandling(err);

	return UTXOs;
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int);
	accumulated := 0;
	db := u.Mattchain.Database;

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions;

		it := txn.NewIterator(opts);
		defer it.Close();

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			item := it.Item();
			k := item.Key();

			var v []byte;

			err := item.Value(func(val []byte) error {
				v = val;
				return nil;
			});
			utils.ErrorHandling(err);

			k = bytes.TrimPrefix(k, utxoPrefix);
			txID := hex.EncodeToString(k);
			outs := DeserializeOutputs(v);

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value;
					unspentOuts[txID] = append(unspentOuts[txID], outIdx);
				}
			}
		}

		return nil;
	})
	utils.ErrorHandling(err);

	return accumulated, unspentOuts;
}

func (u UTXOSet) CountTransactions() int {
	db := u.Mattchain.Database;
	counter := 0;

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions;

		it := txn.NewIterator(opts);
		defer it.Close();

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++;
		}

		return nil;
	});
	utils.ErrorHandling(err);
	
	return counter;
}

func (u UTXOSet) Reindex() {
	db := u.Mattchain.Database;

	u.DeleteByPrefix(utxoPrefix);

	UTXO := u.Mattchain.FindUTXO();

	err := db.Update(func(txn *badger.Txn) error {
		for txId, outs := range UTXO {
			key, err := hex.DecodeString(txId);

			if err != nil {
				return err;
			}

			key = append(utxoPrefix, key...);

			err = txn.Set(key, outs.SerializeOutputs());
			utils.ErrorHandling(err);
		}

		return nil;
	});
	utils.ErrorHandling(err);
}

func (u *UTXOSet) Update(block *MattBlock) {
	db := u.Mattchain.Database;
	
	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					updatedOuts := TxOutputs{};
					inID := append(utxoPrefix, in.ID...);

					item, err := txn.Get(inID);
					utils.ErrorHandling(err);

					var v []byte;

					err = item.Value(func(val []byte) error {
						v = val;
						return nil;
					});
					utils.ErrorHandling(err);

					outs := DeserializeOutputs(v);

					for outIdx, out := range outs.Outputs {
						if outIdx != in.Out {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out);
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							log.Panic(err);
						}
					} else {
						if err := txn.Set(inID, updatedOuts.SerializeOutputs()); err != nil {
							log.Panic(err);
						}
					}
				}
			}

			newOutputs := TxOutputs{};
			newOutputs.Outputs = append(newOutputs.Outputs, tx.Outputs...);

			txID := append(utxoPrefix, tx.ID...);
			
			if err := txn.Set(txID, newOutputs.SerializeOutputs()); err != nil {
				log.Panic(err);
			}
		}

		return nil;
	});

	utils.ErrorHandling(err);
}