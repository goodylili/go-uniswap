package uniswap

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"
	"time"
	"tinkerer/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RouterService struct
type RouterService struct {
	client        *ethclient.Client
	routerAddress common.Address
}

// NewRouterService creates a new RouterService
func NewRouterService(client *ethclient.Client, routerAddress common.Address) *RouterService {
	return &RouterService{
		client:        client,
		routerAddress: routerAddress,
	}
}

//// CalculateMinTokens calculates the minimum tokens to receive based on slippage
//func (s *RouterService) CalculateMinTokens(tokenAddress common.Address, amountEth *big.Int, slippage float64) (*big.Int, error) {
//	estimatedTokens, err := s.GetEstimatedTokensForETH(tokenAddress, amountEth)
//	if err != nil {
//		log.Printf("Error getting estimated tokens: %v", err)
//		return nil, err
//	}
//
//	slippageMultiplier := big.NewFloat(1 - slippage/100)
//	minTokensFloat := new(big.Float).Mul(new(big.Float).SetInt(estimatedTokens), slippageMultiplier)
//
//	minTokens, _ := minTokensFloat.Int(nil)
//
//	log.Printf("Calculated Min Tokens: %s", minTokens)
//
//	return minTokens, nil
//}

// IsValidEthereumAddress checks if an address is a valid Ethereum address
func IsValidEthereumAddress(address string) bool {
	return common.IsHexAddress(address)
}

// CheckUserBalance retrieves the token balance of a user
func (s *RouterService) CheckUserBalance(userAddress string, tokenAddress common.Address) (*big.Int, error) {
	erc20, err := NewERC20(tokenAddress, s.client)
	if err != nil {
		log.Printf("Error creating ERC20 instance: %v", err)
		return nil, err
	}

	address := common.HexToAddress(userAddress)
	balance, err := erc20.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Printf("Error retrieving balance: %v", err)
		return nil, err
	}

	return balance, nil
}

// SwapETHForToken performs a swap from ETH to the specified token
func (s *RouterService) SwapETHForToken(userWalletPrivateKey string, tokenAddress common.Address, amountInEth, minTokens *big.Int) (string, error) {
	cfg := config.LoadConfig()

	router, err := NewUniswapRouter(s.routerAddress, s.client)
	if err != nil {
		log.Printf("Error creating Uniswap V3 router: %v", err)
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(userWalletPrivateKey)
	if err != nil {
		log.Printf("Error parsing private key: %v", err)
		return "", err
	}

	chainID, err := s.client.NetworkID(context.Background())
	if err != nil {
		log.Printf("Error getting network ID: %v", err)
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("Invalid private key")
		return "", errors.New("invalid private key")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("Error getting nonce: %v", err)
		return "", err
	}

	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("Error suggesting gas price: %v", err)
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Printf("Error creating transactor: %v", err)
		return "", err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amountInEth       // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	exactInputSingleParams := ISwapRouterExactInputSingleParams{
		TokenIn:           common.HexToAddress(cfg.WethAddress),
		TokenOut:          tokenAddress,
		Fee:               big.NewInt(3000), // 0.3% pool fee
		Recipient:         fromAddress,
		Deadline:          big.NewInt(time.Now().Unix() + 300),
		AmountIn:          amountInEth,
		AmountOutMinimum:  minTokens,
		SqrtPriceLimitX96: big.NewInt(0),
	}

	tx, err := router.ExactInputSingle(auth, exactInputSingleParams)
	if err != nil {
		log.Printf("Error executing swap: %v", err)
		return "", err
	}

	transactionFee := new(big.Int).Mul(amountInEth, big.NewInt(1))
	transactionFee.Div(transactionFee, big.NewInt(100))

	botWalletAddress := common.HexToAddress(os.Getenv("BOT_WALLET_ADDRESS"))
	feeTx := types.NewTransaction(nonce+1, botWalletAddress, transactionFee, uint64(30000), gasPrice, nil)

	signedFeeTx, err := types.SignTx(feeTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Printf("Error signing transaction: %v", err)
		return "Error signing transaction", err
	}

	err = s.client.SendTransaction(context.Background(), signedFeeTx)
	if err != nil {
		log.Printf("Error sending transaction fee: %v", err)
		return "Error sending transaction fee", err
	}

	log.Printf("Swap executed, transaction hash: %s\n", tx.Hash().Hex())

	return tx.Hash().Hex(), nil
}

//// GetEstimatedTokensForETH estimates the number of tokens that can be received for a given amount of ETH
//func (s *RouterService) GetEstimatedTokensForETH(tokenAddress common.Address, amountEth *big.Int) (*big.Int, error) {
//	router, err := NewUniswapRouter(s.routerAddress, s.client)
//	if err != nil {
//		log.Printf("Error creating Uniswap V3 router: %v", err)
//		return nil, err
//	}
//
//	ethNativeTokenAddress := common.HexToAddress(os.Getenv("WETH_ADDRESS"))
//
//	callOpts := &bind.CallOpts{
//		Pending: false,
//		From:    ethNativeTokenAddress,
//		Context: context.Background(),
//	}
//
//	// Add the appropriate function to estimate the token amount using Uniswap V3 here
//	// For simplicity, this function can be adjusted to fit the actual implementation needed
//
//	// Example placeholder implementation
//	amountsOut, err := router.ExactInputSingle(callOpts, ISwapRouterExactInputSingleParams{
//		TokenIn:          ethNativeTokenAddress,
//		TokenOut:         tokenAddress,
//		Fee:              big.NewInt(3000), // 0.3% pool fee
//		AmountIn:         amountEth,
//		AmountOutMinimum: big.NewInt(1),
//	})
//	if err != nil {
//		log.Printf("Error getting estimated tokens: %v", err)
//		return nil, err
//	}
//
//	log.Printf("Estimated Tokens for %s: %s", tokenAddress.Hex(), amountsOut)
//
//	return amountsOut, nil
//}
