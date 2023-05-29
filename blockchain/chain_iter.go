package blockchain

import (
	"github.com/dgraph-io/badger"
	"github.com/dustinggreenery/MattChain/utils"
)

type MattChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (iter *MattChainIterator) Next() *MattBlock {
	var block *MattBlock;

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash);
		utils.ErrorHandling(err);

		var encodedBlock []byte;

		err = item.Value(func(val []byte) error {
			encodedBlock = val;
			return nil;
		});
		block = DeserializeBlock(encodedBlock);

		return err;
	});
	utils.ErrorHandling(err);

	iter.CurrentHash = block.PrevHash;

	return block;
}