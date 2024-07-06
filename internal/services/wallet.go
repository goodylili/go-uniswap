package services

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// CreateWallet generates a new Ethereum wallet for user.
func CreateWallet() (address, privateKey string, err error) {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	privateKey = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address = string(crypto.PubkeyToAddress(*publicKeyECDSA).Hex())

	return address, privateKey, nil
}

// ImportWallet import an exising Etherum wallet using a private key, from the user.
func ImportWallet(privateKey string) (address string, err error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address = string(crypto.PubkeyToAddress(*publicKeyECDSA).Hex())

	return address, nil
}

func GenerateDummyPrivateKeyForTest() {
	privateKey, err := crypto.GenerateKey()

	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
    privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]

    log.Println("Generated Private Key:", privateKeyHex)
}

