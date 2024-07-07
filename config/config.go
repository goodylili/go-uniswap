package config

type NetworkConfig struct {
	// TODO: RPC URLs are set to default values. Replace them with your own. Do not use in production
	RPCURL         string
	RouterAddress  string
	WethAddress    string
	FactoryAddress string
}

type Config struct {
	Networks map[string]NetworkConfig
	// you can add more configurations here
}

func LoadConfig() *Config {
	return &Config{
		Networks: map[string]NetworkConfig{
			// TODO: set your own RPC URLs, Router Addresses, and WETH Addresses here
			"mainnet": {
				RPCURL:         "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID",
				RouterAddress:  "0xUniswapV3RouterAddressOnMainnet",
				FactoryAddress: "0xUniswapV3FactoryAddressOnMainnet",
				WethAddress:    "0xWETHAddressOnMainnet",
			},
			"arbitrum": {
				RPCURL:         "https://arb1.arbitrum.io/rpc",
				RouterAddress:  "0xUniswapV3RouterAddressOnArbitrum",
				FactoryAddress: "0xUniswapV3FactoryAddressOnArbitrum",
				WethAddress:    "0xWETHAddressOnArbitrum",
			},
			"base": {
				RPCURL:         "https://base-rpc-url",
				RouterAddress:  "0xUniswapV3RouterAddressOnBase",
				WethAddress:    "0xWETHAddressOnBase",
				FactoryAddress: "0xUniswapV3FactoryAddressOnBase",
			},

			// add more networks deployed on Uniswap v2 here e.g ZkSync
		},
	}
}
