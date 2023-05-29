package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/dustinggreenery/MattChain/utils"
	"github.com/dustinggreenery/MattChain/wallet"
)

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey);

	return bytes.Equal(lockingHash, pubKeyHash);
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil};
	txo.Lock([]byte(address));

	return txo;
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address);
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4];
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash);
}

func (outs TxOutputs) SerializeOutputs() []byte {
	var buffer bytes.Buffer;

	encoder := gob.NewEncoder(&buffer);
	err := encoder.Encode(outs);
	utils.ErrorHandling(err);

	return buffer.Bytes();
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs;

	decoder := gob.NewDecoder(bytes.NewReader(data));
	err := decoder.Decode(&outputs);
	utils.ErrorHandling(err);
	
	return outputs;
}