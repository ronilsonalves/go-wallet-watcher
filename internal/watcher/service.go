package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/ronilsonalves/go-wallet-watcher/internal/domain"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// StartWatcherService load from environment the data and start running goroutines to perform wallet watcher service.
func StartWatcherService() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err.Error())
	}

	var wfe [20]domain.Wallet
	var wallets []domain.Wallet

	for index := range wfe {
		wallet := domain.Wallet{
			Address:   os.Getenv("WATCHER_WALLET" + strconv.Itoa(index+1)),
			SecretKey: os.Getenv("WATCHER_SECRET" + strconv.Itoa(index+1)),
		}
		wallets = append(wallets, wallet)
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(len(wallets))

	// Start a goroutine for each wallet
	for _, wallet := range wallets {
		go func(wallet domain.Wallet) {
			logUrl := os.Getenv("WATCHER_LOGS_URL")
			cType := "application/json"
			authToken := os.Getenv("WATCHER_SOURCE_TOKEN")
			httpCli := &http.Client{}
			// Connect to the Ethereum client
			client, err := rpc.Dial(os.Getenv("WATCHER_RPC_ADDRESS"))
			if err != nil {
				body, _ := json.Marshal(domain.LogMsg{
					Dt:      time.Now().String(),
					Message: "Failed to connect to the RPC client for address " + fmt.Sprintf("%s", wallet.Address),
				})
				req, _ := http.NewRequest("POST", logUrl, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", cType)
				req.Header.Set("Authorization", authToken)
				resp, _ := httpCli.Do(req)
				defer resp.Body.Close()
				//log.Printf("Failed to connect to the RPC client for address %s: %v \n Trying fallback rpc server", wallet.Address.Hex(), err)
			}
			client, err = rpc.Dial(os.Getenv("WATCHER_RPC_FALLBACK_ADDRESS"))
			if err != nil {
				body, _ := json.Marshal(domain.LogMsg{
					Dt:      time.Now().String(),
					Message: "Failed to connect to the RPC client for address " + fmt.Sprintf("%s", wallet.Address),
				})
				req, _ := http.NewRequest("POST", logUrl, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", cType)
				req.Header.Set("Authorization", authToken)
				resp, _ := httpCli.Do(req)
				defer resp.Body.Close()
				//log.Printf("Failed to connect to the Ethereum client for address %s: %v", wallet.Address.Hex(), err)
				wg.Done()
				return
			}

			// Create an instance of the Ethereum client
			ethClient := ethclient.NewClient(client)

			for {
				// Get the balance of the address
				balance, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(wallet.Address), nil)
				if err != nil {
					body, _ := json.Marshal(domain.LogMsg{
						Dt:      time.Now().String(),
						Message: "Failed to get balance for address " + fmt.Sprintf("%s: ", wallet.Address),
					})
					req, _ := http.NewRequest("POST", logUrl, bytes.NewBuffer(body))
					req.Header.Set("Content-Type", cType)
					req.Header.Set("Authorization", authToken)
					resp, _ := httpCli.Do(req)
					resp.Body.Close()
					//defer resp.Body.Close()
					//log.Printf("Failed to get balance for address %s: %v", wallet.Address.Hex(), err)
					continue
				}

				balanceInEther := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

				//log.Printf("Balance for address %s: %.16f ETH", wallet.Address.Hex(), balanceInEther)
				body, _ := json.Marshal(domain.LogMsg{
					Dt: time.Now().String(),
					Message: "Balance for address " + fmt.Sprintf("%s", wallet.Address) + ": " +
						fmt.Sprintf("%.18f ETH", balanceInEther),
				})
				req, _ := http.NewRequest("POST", logUrl, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", cType)
				req.Header.Set("Authorization", authToken)
				resp, _ := httpCli.Do(req)
				resp.Body.Close()

				// if the wallet has a balance superior to 0.0005 ETH, we are sending the balance to another wallet
				if balanceInEther.Cmp(big.NewFloat(0.0005)) > 0 {
					sendBalanceToAnotherWallet(common.HexToAddress(wallet.Address), balance, wallet.SecretKey)
				}

				time.Sleep(300 * time.Millisecond) // Wait for a while before checking for the next block
			}
		}(wallet)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// sendBalanceToAnotherWallet when find some values in any wallet perform a SendTransaction(ctx context.Context,
// tx *types.Transaction) function
func sendBalanceToAnotherWallet(fromAddress common.Address, balance *big.Int, privateKeyHex string) {
	toAddress := common.HexToAddress(os.Getenv("WATCHER_DEST_ADDRESS"))
	chainID := big.NewInt(1) // Replace with the appropriate chain ID

	// Connect to the Ethereum client
	client, err := rpc.Dial(os.Getenv("WATCHER_RPC_ADDRESS"))
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v...", err)
	}

	ethClient := ethclient.NewClient(client)

	// Load the private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Get the current nonce for the fromAddress
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("Failed to retrieve nonce: %v", err)
	}

	// Create a new transaction
	gasLimit := uint64(21000) // Set the gas limit based on the transaction type
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("Failed to retrieve gas price: %v", err)
	}

	//tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    new(big.Int).Sub(balance, new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))),
		Data:     nil,
	})
	valueInEther := new(big.Float).Quo(new(big.Float).SetInt(tx.Value()), big.NewFloat(1e18))
	if valueInEther.Cmp(big.NewFloat(0)) < 0 {
		log.Println("ERROR: Insufficient funds to make transfer")
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Printf("Failed to sign transaction: %v", err)
	}

	// Send the signed transaction
	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
	} else {
		log.Printf("Transaction sent: %s", signedTx.Hash().Hex())
	}
}
