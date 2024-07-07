## Overview

The `config` package helps you manage network configurations for different blockchain networks, such as RPC URLs and contract addresses.

You'll work with the main structures `NetworkConfig` and `Config`. The `LoadConfig` function also initializes and returns a configuration instance with default settings.

### `NetworkConfig`

The `NetworkConfig` struct holds configuration details for a specific network, including the RPC URL and various contract addresses.

```go

type NetworkConfig struct {

    RPCURL         string

    RouterAddress  string

    WethAddress    string

    FactoryAddress string

}

```

### `Config`

The `Config` struct contains a map of network configurations. You can expand it by adding more configurations as needed.

```go

type Config struct {

    Networks map[string]NetworkConfig

}

```

### `LoadConfig` Function

The `LoadConfig` function initializes the configuration with default values. Remember to replace these defaults with your own specific values before using it in a production environment.

```go

func LoadConfig() *Config {

    return &Config{

        Networks: map[string]NetworkConfig{

            "mainnet": {

                RPCURL:        "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID",

                RouterAddress: "0xUniswapV3RouterAddressOnMainnet",

                FactoryAddress: "0xUniswapV3FactoryAddressOnMainnet",

                WethAddress:   "0xWETHAddressOnMainnet",

            },

            "arbitrum": {

                RPCURL:        "https://arb1.arbitrum.io/rpc",

                RouterAddress: "0xUniswapV3RouterAddressOnArbitrum",

                FactoryAddress: "0xUniswapV3FactoryAddressOnArbitrum",

                WethAddress:   "0xWETHAddressOnArbitrum",

            },

            "base": {

                RPCURL:        "https://base-rpc-url",

                RouterAddress: "0xUniswapV3RouterAddressOnBase",

                FactoryAddress: "0xUniswapV3FactoryAddressOnBase",

                WethAddress:   "0xWETHAddressOnBase",

            },

            // Add more networks deployed on Uniswap v2 here, e.g., ZkSync

        },

    }

}

```

## How to Use

1. Import the Package:

   Make sure to import the `config` package in your Go file.

   ```go

   import "your_project/config"

   ```

2. Load the Configuration:

   Call the `LoadConfig` function to initialize the configuration with default values.

   ```go

   cfg := config.LoadConfig()

   ```

3. Customize the Configuration:

   Replace the default RPC URLs, Router Addresses, and WETH Addresses with your specific values.

   ```go

   cfg.Networks["mainnet"] = config.NetworkConfig{

       RPCURL:        "https://mainnet.infura.io/v3/YOUR_NEW_INFURA_PROJECT_ID",

       RouterAddress: "0xYourNewUniswapV3RouterAddressOnMainnet",

       FactoryAddress: "0xYourNewUniswapV3FactoryAddressOnMainnet",

       WethAddress:   "0xYourNewWETHAddressOnMainnet",

   }

   ```

4. Add More Networks:

   You can add more network configurations to the `Networks` map by following the structure provided.

   ```go

   cfg.Networks["zksync"] = config.NetworkConfig{

       RPCURL:        "https://zksync-rpc-url",

       RouterAddress: "0xUniswapV3RouterAddressOnZkSync",

       FactoryAddress: "0xUniswapV3FactoryAddressOnZkSync",

       WethAddress:   "0xWETHAddressOnZkSync",

   }

   ```

## Notes

- Security: Ensure that the RPC URLs and contract addresses are kept secure and are not exposed in your source code, especially in production environments.

- Customization: Feel free to extend the `Config` struct with additional configurations as needed.

That's it! Using the ' config ' package, you've successfully set up and customized your network configurations. 