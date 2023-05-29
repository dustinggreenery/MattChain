package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/dustinggreenery/MattChain/merkle"
	"github.com/dustinggreenery/MattChain/utils"
)

type MattBlock struct {
	Hash         []byte
	Transactions []*MattTransaction
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(txs []*MattTransaction, prevHash []byte) *MattBlock {
	block := &MattBlock{[]byte{}, txs, prevHash, 0};
	
	pow := NewProof(block);
	nonce, hash := pow.Run();

	block.Hash = hash[:];
	block.Nonce = nonce;

	return block;
}

func Genesis(coinbase *MattTransaction) *MattBlock {
	return CreateBlock([]*MattTransaction{coinbase}, []byte{});
}

func (b *MattBlock) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}

	tree := merkle.NewMerkleTree(txHashes);

	return tree.RootNode.Data;
}

func (b *MattBlock) SerializeBlock() []byte {
	var res bytes.Buffer;
	encoder := gob.NewEncoder(&res);

	err := encoder.Encode(b);

	utils.ErrorHandling(err);

	return res.Bytes();
}

func DeserializeBlock(data []byte) *MattBlock {
	var block MattBlock;

	decoder := gob.NewDecoder(bytes.NewReader(data));
	err := decoder.Decode(&block);
	utils.ErrorHandling(err);

	return &block;
}