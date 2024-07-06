package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"tinkerer/internal/uniswap"
)

type TradeService struct {
	client        *ethclient.Client
	routerService *uniswap.RouterService
}

func InitializeUniswapContracts() (*ethclient.Client, common.Address, error) {
	WEB3_PROVIDER_URL := os.Getenv("WEB3_PROVIDER_URL")
	if WEB3_PROVIDER_URL == "" {
		log.Println("WEB3_PROVIDER_URL is not set.")
	}

	// Initialize Ethereum client
	client, err := ethclient.Dial(WEB3_PROVIDER_URL)
	if err != nil {
		log.Printf("Error dialing Ethereum client: %v\n", err)
		return nil, common.Address{}, err
	}

	// Router contract address
	routerAddress := common.HexToAddress("0xf164fC0Ec4E93095b804a4795bBe1e041497b92a")

	return client, routerAddress, nil
}

func NewTradeService() (*TradeService, error) {
	client, routerAddress, err := InitializeUniswapContracts()
	if err != nil {
		return nil, err
	}

	routerService := uniswap.NewRouterService(client, routerAddress)

	return &TradeService{
		client:        client,
		routerService: routerService,
	}, nil
}

func (s *TradeService) ExecuteTrade(userWalletPrivateKey, tokenAddress string, amountEth *big.Int, slippage float64) (string, error) {
	// Validate the private key.
	if _, err := crypto.HexToECDSA(userWalletPrivateKey); err != nil {
		log.Printf("Invalid private key: %v", err)
		return "", errors.New("invalid private key")
	}

	// Convert the token address to common.Address
	tokenAddr := common.HexToAddress(tokenAddress)
	if !uniswap.IsValidEthereumAddress(tokenAddr.Hex()) {
		return "", errors.New("invalid token address")
	}

	// Convert the private key from hex to *ecdsa.PrivateKey
	privateKey, err := crypto.HexToECDSA(userWalletPrivateKey)
	if err != nil {
		log.Printf("Invalid private key: %v", err)
		return "", errors.New("invalid private key")
	}

	// Derive the public key from the private key
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Println("Error casting public key to ECDSA")
		return "", errors.New("error casting public key to ECDSA")
	}

	// Derive the Ethereum address from the public key.
	userAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Check user's ETH balance
	ethBalance, err := s.client.BalanceAt(context.Background(), userAddress, nil)
	if err != nil {
		log.Printf("Error retrieving ETH balance: %v", err)
		return "", err
	}

	if ethBalance.Cmp(amountEth) < 0 {
		return "", errors.New("insufficient ETH balance for trade")
	}

	// Calculate min tokens to accept based on slippage.
	minTokens, err := s.routerService.CalculateMinTokens(tokenAddr, amountEth, slippage)
	if err != nil {
		log.Fatalf("Failed to calculate min tokens, due to: %v", err)
		return "Failed to calculate minimum tokens based on slippage", err
	}

	// Execute the swap.
	txHash, err := s.routerService.SwapETHForToken(userWalletPrivateKey, tokenAddr, amountEth, minTokens)
	if err != nil {
		log.Fatalf("Failed to swap ETH for the token, due to: %v", err)
		return "Failed to swap ETH for the token", err
	}

	return txHash, nil
}
