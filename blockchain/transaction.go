package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/dustinggreenery/MattChain/utils"
	"github.com/dustinggreenery/MattChain/wallet"
)

type MattTransaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx MattTransaction) Serialize() []byte {
	var encoded bytes.Buffer;
	
	enc := gob.NewEncoder(&encoded);
	
	err := enc.Encode(tx);
	utils.ErrorHandling(err);

	return encoded.Bytes();
}

func (tx *MattTransaction) Hash() []byte {
	var hash [32]byte;

	txCopy := *tx;
	txCopy.ID = []byte{};

	hash = sha256.Sum256(txCopy.Serialize());

	return hash[:];
}

func (tx *MattTransaction) TrimmedCopy() MattTransaction {
	var inputs []TxInput;
	var outputs []TxOutput;

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil});
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash});
	}

	txCopy := MattTransaction{tx.ID, inputs, outputs};

	return txCopy;
}

func CoinbaseTx(to, data string) *MattTransaction {
	if data == "" {
		randData := make([]byte, 24);

		_, err := rand.Read(randData);
		utils.ErrorHandling(err);

		data = fmt.Sprintf("%x", randData);
	}

	txin := TxInput{[]byte{}, -1, nil, []byte(data)};
	txout := NewTXOutput(20, to);

	tx := MattTransaction{nil, []TxInput{txin}, []TxOutput{*txout}};
	tx.ID = tx.Hash();

	return &tx;
}

func NewTransaction(from, to string, amount int, UTXO *UTXOSet) *MattTransaction {
	var inputs []TxInput;
	var outputs []TxOutput;

	wallets, err := wallet.CreateWallets();
	utils.ErrorHandling(err);

	w := wallets.GetWallet(from);
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey);
	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount);

	if acc < amount {
		log.Panic("Error: Not Enough Funds to Complete Transaction");
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid);
		utils.ErrorHandling(err);

		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey};
			inputs = append(inputs, input);
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to));

	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from));
	}

	tx := MattTransaction{nil, inputs, outputs};
	tx.ID = tx.Hash();
	UTXO.Mattchain.SignTransaction(&tx, w.PrivateKey);

	return &tx;
}

func (tx *MattTransaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 &&
		   len(tx.Inputs[0].ID) == 0 && 
		   tx.Inputs[0].Out == -1;
}

func (tx *MattTransaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]MattTransaction) {
	if tx.IsCoinbase() {
		return;
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Error: Previous transaction doesn't exist");
		}
	}

	txCopy := tx.TrimmedCopy();

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)];

		txCopy.Inputs[inId].Signature = nil;
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash;

		dataToSign := fmt.Sprintf("%x\n", txCopy);

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign));
		utils.ErrorHandling(err);

		signature := append(r.Bytes(), s.Bytes()...);

		tx.Inputs[inId].Signature = signature;
		txCopy.Inputs[inId].PubKey = nil;
	}
}

func (tx *MattTransaction) Verify(prevTXs map[string]MattTransaction) bool {
	if tx.IsCoinbase() {
		return true;
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("Previous transaction does not exist");
		}
	}

	txCopy := tx.TrimmedCopy();
	curve := elliptic.P256();

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)];
		txCopy.Inputs[inId].Signature = nil;
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash;

		r := big.Int{};
		s := big.Int{};
		sigLen := len(in.Signature);
		r.SetBytes(in.Signature[:(sigLen / 2)]);
		s.SetBytes(in.Signature[(sigLen / 2):]);

		x := big.Int{};
		y := big.Int{};
		keyLen := len(in.PubKey);
		x.SetBytes(in.PubKey[:(keyLen / 2)]);
		y.SetBytes(in.PubKey[(keyLen / 2):]);

		dataToVerify := fmt.Sprintf("%x\n", txCopy);

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y};

		if !ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) {
			return false;
		}

		txCopy.Inputs[inId].PubKey = nil;
	}

	return true;
}

func (tx MattTransaction) String() string {
	var lines []string;

	lines = append(lines, fmt.Sprintf("  Transaction %x:", tx.ID));

	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("    Input %d:", i));
		lines = append(lines, fmt.Sprintf("      TXID:	 %x", input.ID));
		lines = append(lines, fmt.Sprintf("      Out:	 %d", input.Out));
		lines = append(lines, fmt.Sprintf("      Signature: %x", input.Signature));
		lines = append(lines, fmt.Sprintf("      PubKey:	 %x", input.PubKey));
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("    Output %d:", i));
		lines = append(lines, fmt.Sprintf("      Value:	 %d", output.Value));
		lines = append(lines, fmt.Sprintf("      Script:	 %x", output.PubKeyHash));
	}

	return strings.Join(lines, "\n");
}