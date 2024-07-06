package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
)

var (
	logger, _ = zap.NewDevelopment()
)

type Config struct {
	Web3ProviderURL       string
	WethAddress           string
	UniswapRouterAddress  string
	UniswapFactoryAddress string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	return &Config{
		Web3ProviderURL:       os.Getenv("WEB3_PROVIDER_URL"),
		WethAddress:           os.Getenv("WETH_ADDRESS"),
		UniswapRouterAddress:  os.Getenv("UNISWAP_ROUTER_ADDRESS"),
		UniswapFactoryAddress: os.Getenv("UNISWAP_FACTORY_ADDRESS"),
	}
}
