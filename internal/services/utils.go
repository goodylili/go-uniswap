package services

//
//import (
//	"fmt"
//	"log"
//	"math/big"
//	"strings"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/ethereum/go-ethereum/ethclient"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"/internal/database"
//	"github.com/theghostmac/tugbot/internal/uniswap"
//)
//
//// Helper function to get the verification status emoji
//func getVerificationStatusEmoji(isVerified bool) string {
//	if isVerified {
//		return "✅"
//	}
//	return "❌"
//}
//
//func HandleInlineQuery(bot *tgbotapi.BotAPI, callBackQuery *tgbotapi.CallbackQuery) {
//	// Splitting the callback data to understand the action.
//	parts := strings.Split(callBackQuery.Data, ":")
//	action := parts[0]
//
//	switch action {
//	case "trade":
//		// New logic to handle trade action.
//		if len(parts) != 3 {
//			log.Printf("Invalid callback data for trade: %s", callBackQuery.Data)
//			return
//		}
//		tokenAddress := parts[1]
//		amountEthString := parts[2]
//
//		// Convert ETH amount to Wei (1 ETH = 1e18 Wei)
//		amountEthInWei := new(big.Float).Mul(big.NewFloat(1e18), big.NewFloat(0))
//		_, ok := amountEthInWei.SetString(amountEthString)
//		if !ok {
//			HandleBotMessages(bot, callBackQuery.Message.Chat.ID, "Invalid ETH amount format")
//			return
//		}
//
//		// Convert amount from big.Float to big.Int
//		amountEth := new(big.Int)
//		amountEthInWei.Int(amountEth) // Converts and stores the result in amountEth
//
//
//		// Fetch user's stored private key
//		privateKey, err := database.GetUserPrivateKey(callBackQuery.From.ID)
//		// User does not have a private key, redirect to wallet setup.
//		if err != nil || privateKey == "" {
//			redirectToPrivateChat(bot, callBackQuery.From.ID, tokenAddress)
//			database.SetUserState(callBackQuery.From.ID, fmt.Sprintf("awaitingWalletSetup:%s:%s", tokenAddress, amountEthString))
//		} else {
//			// User has a private key, execute the trade.
//			executeTrade(bot, callBackQuery, privateKey, tokenAddress, amountEth)
//		}
//	case "create_wallet":
//		// The user selected to create a new wallet.
//		handleCreateWallet(bot, callBackQuery)
//
//	case "import_wallet":
//		// The user selected to import an existing wallet.
//		database.SetUserState(callBackQuery.Message.Chat.ID, "awaitingPrivateKey")
//		msg := tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Please enter your private key to import your wallet.")
//		bot.Send(msg)
//
//	case "setup_wallet":
//		// WalletCommand(bot, callBackQuery.Message)
//		// // TODO: delete this:
//		// Handle wallet setup action.
//		handleWalletSetupInline(bot, callBackQuery)
//
//	case "disconnect_wallet":
//		handleDisconnectWallet(bot, callBackQuery)
//	case "check_balance":
//		handleCheckBalance(bot, callBackQuery)
//	}
//}
//
//func executeTrade(bot *tgbotapi.BotAPI, callBackQuery *tgbotapi.CallbackQuery, privateKey string, tokenAddress string, amountEth *big.Int) {
//	// Initialize the trade service.
//	tradeService, err := NewTradeService()
//	if err != nil {
//		sendMessageToUser(bot, callBackQuery.From.ID, "Error initializing trade service: "+err.Error())
//		return
//	}
//
//	// Set a default slippage of 0.5%
//	txHash, err := tradeService.ExecuteTrade(privateKey, tokenAddress, amountEth, 0.005)
//	if err != nil {
//		sendMessageToUser(bot, callBackQuery.Message.Chat.ID, "Please check the DM for the trade details.")
//		sendMessageToUser(bot, callBackQuery.From.ID, "Trade execution failed: "+err.Error())
//		return
//	}
//
//	// Respond with the transaction hash.
//	response := "Trade executed successfully. Transaction Hash: " + txHash
//	sendMessageToUser(bot, callBackQuery.From.ID, response)
//}
//
//// handleWalletSetupCompletion continues the trade after wallet creation or import in the wallet setup function.
//func handleWalletSetupCompletion(bot *tgbotapi.BotAPI, userID int64) {
//    log.Println("handleWalletSetupCompletion triggered")  // Log entry for function trigger
//
//    userState, err := database.GetUserState(userID)
//    if err != nil {
//        log.Printf("Error fetching user state: %v", err)
//        return
//    }
//
//    // Log the current user state
//    log.Printf("Current user state: %s", userState)
//
//    // Check if the user just finished setting up the wallet.
//    if strings.HasPrefix(userState, "awaitingWalletSetup") {
//        parts := strings.Split(userState, ":")
//        if len(parts) == 3 {
//            tokenAddress := parts[1]
//            amountEthString := parts[2]
//
//            // Update state to awaitingTrade with tokenAddress and amountEthString
//            newState := fmt.Sprintf("awaitingTrade:%s:%s", tokenAddress, amountEthString)
//            if err := database.SetUserState(userID, newState); err != nil {
//                log.Printf("Error setting user state to awaitingTrade: %v", err)
//            } else {
//                log.Printf("User state updated to: %s", newState)  // Log the new state
//            }
//
//            // Optionally, send a confirmation message to the user
//            sendMessageToUser(bot, userID, "Wallet setup completed. Proceeding with the trade.")
//        }
//    }
//
//    // Check if the user is ready to execute the trade.
//    if strings.HasPrefix(userState, "awaitingTrade") {
//        parts := strings.Split(userState, ":")
//        if len(parts) == 3 {
//            tokenAddress := parts[1]
//            amountEthString := parts[2]
//
//            // Convert the amount from string to big.Int
//            amountEth, ok := new(big.Int).SetString(amountEthString, 10)
//            if !ok {
//                sendMessageToUser(bot, userID, "Invalid amount format in trade intent.")
//                return
//            }
//
//            // Fetch the stored private key.
//            privateKey, err := database.GetUserPrivateKey(userID)
//            if err != nil || privateKey == "" {
//                log.Printf("Error fetching private key: %v", err)
//                sendMessageToUser(bot, userID, "Error fetching private key.")
//                return
//            }
//
//            log.Printf("Fetched private key for user %d: %s", userID, privateKey) // Log the fetched private key
//
//            // Execute the trade in DM.
//            executeTradeInDM(bot, userID, privateKey, tokenAddress, amountEth)
//        }
//
//        // Clear the user state after handling the trade intent.
//        database.ClearUserState(userID)
//    }
//
//	if userState == "walletSetupComplete" {
//		// proceed with the trade.
//		initiateTrade(bot, userID)
//	}
//}
//
//func sendMessageToUser(bot *tgbotapi.BotAPI, userID int64, messageText string) {
//    msg := tgbotapi.NewMessage(userID, messageText)
//    _, err := bot.Send(msg)
//    if err != nil {
//        log.Printf("Error sending message to user %d: %v", userID, err)
//    }
//}
//
//func executeTradeInDM(bot *tgbotapi.BotAPI, userID int64, privateKey, tokenAddress string, amountEth *big.Int) {
//	// Inform user the trade is being processed.
//	sendMessageToUser(bot, userID, "Processing your trade...")
//	// Initialize trade service.
//	tradeService, err := NewTradeService()
//	if err != nil {
//		sendMessageToUser(bot, userID, "Error initializing trade service: "+err.Error())
//		return
//	}
//
//	// Execute the trade with a default slippage of 0.5%
//	txHash, err := tradeService.ExecuteTrade(privateKey, tokenAddress, amountEth, 0.005)
//	if err != nil {
//		sendMessageToUser(bot, userID, "Trade execution failed: "+err.Error())
//		return
//	}
//
//	// Send the transaction hash to the user.
//	sendMessageToUser(bot, userID, "Trade executed successfully. Transaction Hash: "+txHash)
//}
//
//// handleDisconnectWallet removes user's private key from the database.
//func handleDisconnectWallet(bot *tgbotapi.BotAPI, callBackQuery *tgbotapi.CallbackQuery) {
//	err := database.RemoveUser(callBackQuery.From.ID)
//	if err != nil {
//		msg := tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Failed to disconnect wallet: "+err.Error())
//		bot.Send(msg)
//		return
//	}
//
//	msg := tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Wallet disconnected successfully.")
//	bot.Send(msg)
//}
//
//// handleCheckBalance checks the balance of the user's wallet.
//func handleCheckBalance(bot *tgbotapi.BotAPI, callBackQuery *tgbotapi.CallbackQuery) {
//	privateKey, err := database.GetUserPrivateKey(callBackQuery.From.ID)
//	if err != nil || privateKey == "" {
//		msg := tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "No wallet connected.")
//		bot.Send(msg)
//		return
//	}
//
//	// Convert private key to public address.
//	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
//	if err != nil {
//		bot.Send(tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Failed to parse private key"))
//		return
//	}
//
//	publicAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
//
//	// Connect to the Etherum client.
//	client, err := ethclient.Dial("WEB3_PROVIDER_URL")
//	if err != nil {
//		bot.Send(tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Failed to connect to Ethereum client"))
//		return
//	}
//	defer client.Close()
//
//	// Create a new instance of the ERC20 contract.
//	tokenAddress := common.HexToAddress("WETH_ADDRESS")
//	instance, err := uniswap.NewERC20(tokenAddress, client)
//	if err != nil {
//		bot.Send(tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Failed to create ERC20 instance"))
//		return
//	}
//
//	// Retrieve the balance.
//	balance, err := instance.BalanceOf(nil, publicAddress)
//	if err != nil {
//		bot.Send(tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, "Failed to retrieve balance"))
//		return
//	}
//
//	// Send the balance to the user.
//	msg := fmt.Sprintf("Your balance is: %s", balance.String())
//	bot.Send(tgbotapi.NewMessage(callBackQuery.Message.Chat.ID, msg))
//}
//
//func initiateTrade(bot *tgbotapi.BotAPI, userID int64) {
//    // Fetch the stored private key
//    privateKey, err := database.GetUserPrivateKey(userID)
//    if err != nil || privateKey == "" {
//        sendMessageToUser(bot, userID, "Error fetching private key, please set up your wallet.")
//        return
//    }
//
//    // Fetch the token address and amount for the trade from the user's state
//    // Assuming the token address and amount are stored in the user's state
//    // For example, if the state was set as "walletSetupComplete:tokenAddress:amount"
//    tokenAddress, amountEthString, err := getTokenAddressAndAmountFromState(userID)
//    if err != nil {
//        sendMessageToUser(bot, userID, "Error getting trade details: "+err.Error())
//        return
//    }
//
//    // Convert the amount from string to big.Int
//    amountEth, ok := new(big.Int).SetString(amountEthString, 10)
//    if !ok {
//        sendMessageToUser(bot, userID, "Invalid amount format for the trade.")
//        return
//    }
//
//    // Execute the trade
//    executeTradeInDM(bot, userID, privateKey, tokenAddress, amountEth)
//
//    // Clear the user state after handling the trade intent
//    database.ClearUserState(userID)
//}
//
//// getTokenAddressAndAmountFromState retrieves the token address and amount from the user's state.
//func getTokenAddressAndAmountFromState(userID int64) (string, string, error) {
//    userState, err := database.GetUserState(userID)
//    if err != nil {
//        log.Printf("Error fetching user state: %v", err)
//        return "", "", err
//    }
//
//    // Assuming the state format is "walletSetupComplete:tokenAddress:amount"
//    parts := strings.Split(userState, ":")
//    if len(parts) != 3 {
//        return "", "", fmt.Errorf("invalid state format")
//    }
//
//    tokenAddress := parts[1]
//    amountEthString := parts[2]
//
//    return tokenAddress, amountEthString, nil
//}
