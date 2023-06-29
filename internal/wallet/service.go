package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/ronilsonalves/go-wallet-watcher/internal/domain"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

type Service interface {
	GetWalletBalanceByAddress(address string) (domain.Wallet, error)
	GetTransactionsByAddress(address, page, size string) (domain.Wallet, error)
}

type service struct{}

// NewService creates a new instance of the Wallet Service.
func NewService() Service {
	return &service{}
}

// GetWalletBalanceByAddress retrieves the wallet balance for the given address
func (s service) GetWalletBalanceByAddress(address string) (domain.Wallet, error) {
	// Retrieves Etherscan.io API Key from environment
	apiKey := os.Getenv("WATCHER_ETHERSCAN_API")
	url := fmt.Sprintf(fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s", address, apiKey))

	// Send GET request to the Etherscan API
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make Etherscan API request: %v", err)
		return domain.Wallet{}, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return domain.Wallet{}, err
	}

	// Creates a struct to represent etherscan API response
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	// Parse the JSON response
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return domain.Wallet{}, err
	}

	if result.Status != "1" {
		log.Printf("API returned error: %s", result.Message)
		return domain.Wallet{}, fmt.Errorf("API error: %s", result.Message)
	}

	wbBigInt := new(big.Int)
	wbBigInt, ok := wbBigInt.SetString(result.Result, 10)
	if !ok {
		log.Println("Failed to parse string to BigInt")
		return domain.Wallet{}, fmt.Errorf("failed to parse string into BigInt. result.Result value: %s", result.Result)
	}

	wb := new(big.Float).Quo(new(big.Float).SetInt(wbBigInt), big.NewFloat(1e18))
	v, _ := strconv.ParseFloat(wb.String(), 64)

	return domain.Wallet{
		Address: address,
		Balance: v,
	}, nil
}

// GetTransactionsByAddress retrieves the wallet balance and last transactions for the given address paggeable
func (s service) GetTransactionsByAddress(address, page, size string) (domain.Wallet, error) {
	wallet, _ := s.GetWalletBalanceByAddress(address)
	apiKey := os.Getenv("WATCHER_ETHERSCAN_API")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=%s&offset=%s&sort=desc&apikey=%s", address, page, size, apiKey)

	// Send GET request to the Etherscan API
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make API request: %v", err)
		return domain.Wallet{}, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return domain.Wallet{}, err
	}

	// Parse the JSON response
	var transactions struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Result  []domain.Transaction `json:"result"`
	}
	err = json.Unmarshal(body, &transactions)
	if err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return domain.Wallet{}, err
	}

	if transactions.Status != "1" {
		log.Printf("API returned error: %s", transactions.Message)
		return domain.Wallet{}, fmt.Errorf("API error: %s", transactions.Message)
	}

	wallet.Transactions = append(wallet.Transactions, transactions.Result...)

	return wallet, nil
}
