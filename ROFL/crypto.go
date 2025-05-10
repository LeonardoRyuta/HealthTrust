package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GenerateKeyPair creates a secure ECIES key pair
func GenerateKeyPair() (privateKeyHex string, publicKeyHex string, err error) {
	// Generate a secure ECDSA private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Extract public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("failed to cast public key to ECDSA")
	}

	log.Println("Private Key:", privateKey)
	log.Println("Public Key:", publicKeyECDSA)

	// Convert to hex strings for storage
	privateKeyHex = hexutil.Encode(crypto.FromECDSA(privateKey))
	publicKeyHex = hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))
	log.Printf("Private Key: %s", privateKeyHex)
	log.Printf("Public Key: %s", publicKeyHex)

	return privateKeyHex, publicKeyHex, nil
}

func storePubKeyInSC(pubKey string) error {
	log.Printf("Storing public key in SC: %s", pubKey)
	cli, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return err
	}

	// Load the private key
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}

	// Get the account address from the private key
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce for the account
	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get suggested gas price
	gasPrice, err := cli.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Create the ABI
	escAbi, err := abi.JSON(strings.NewReader(ABI_JSON))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Pack the data
	input, err := escAbi.Pack("storePubKey", pubKey)
	if err != nil {
		return fmt.Errorf("failed to pack data: %v", err)
	}

	log.Printf("Storing public key in SC: %s", pubKey)

	// Estimate gas limit
	gasLimit, err := cli.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &CONTRACT_ADDR,
		Data: input,
	})
	if err != nil {
		return fmt.Errorf("failed to estimate gas: %v", err)
	}
	// Add a buffer to the gas limit
	gasLimit = gasLimit * 120 / 100

	// Create the transaction
	tx := types.NewTransaction(
		nonce,
		CONTRACT_ADDR,
		big.NewInt(0), // No ETH being sent
		gasLimit,
		gasPrice,
		input,
	)

	// Get chainID
	chainID, err := cli.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	// Wait for transaction receipt
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), cli, signedTx)
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func DecryptData(encryptedData []byte) (string, error) {
	if privKey == nil {
		return "", errors.New("private key not initialized")
	}

	// Convert hex string private key to ECDSA private key (removing 0x prefix if present)
	privateKeyStr := strings.TrimPrefix(*privKey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// Convert ECDSA private key to ECIES private key
	eciesPrivateKey := ecies.ImportECDSA(privateKey)

	// Handle both hex-encoded and raw byte inputs
	var dataToDecrypt []byte
	// Check if the data starts with 0x (hex format from TypeScript)
	if len(encryptedData) > 2 && string(encryptedData[:2]) == "0x" {
		// Decode from hex if it starts with 0x
		dataToDecrypt, err = hexutil.Decode(string(encryptedData))
		if err != nil {
			return "", fmt.Errorf("failed to decode hex data: %v", err)
		}
	} else {
		// Use raw bytes
		dataToDecrypt = encryptedData
	}

	// Decrypt the data
	plaintext, err := eciesPrivateKey.Decrypt(dataToDecrypt, nil, nil)
	if err != nil {
		// If direct decryption fails, try alternative ECIES format
		// Some ECIES implementations have different byte formats
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	// Return the decrypted text
	return string(plaintext), nil
}
