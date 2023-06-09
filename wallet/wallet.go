package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/dustinggreenery/MattChain/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version = byte(0x00)
)

type MattWallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func (w MattWallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey);
	
	versionedHash := append([]byte{version}, pubHash...);
	checksum := Checksum(versionedHash);

	fullHash := append(versionedHash, checksum...);
	address := utils.Base58Encode(fullHash);

	return address;
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256();

	private, err := ecdsa.GenerateKey(curve, rand.Reader);
	utils.ErrorHandling(err);

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...);
	return *private, pub;
}

func MakeWallet() *MattWallet {
	private, public := NewKeyPair();
	wallet := MattWallet{private, public};

	return &wallet;
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey);

	hasher := ripemd160.New();
	
	_, err := hasher.Write(pubHash[:]);
	utils.ErrorHandling(err);

	publicRipMD := hasher.Sum(nil);

	return publicRipMD;
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload);
	secondHash := sha256.Sum256(firstHash[:]);

	return secondHash[:checksumLength];
}

func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address));
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:];
	version := pubKeyHash[0];
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - checksumLength];
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...));

	return bytes.Equal(actualChecksum, targetChecksum);
}