package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/amargc/simple-go-grpc-wallet/src/wallet"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	client := wallet.NewWalletClient(conn)

	wg.Add(4)

	go performWalletActions(client, "user1", wallet.WALLET_CURRENCY_USD)
	go performWalletActions(client, "user2", wallet.WALLET_CURRENCY_USD)
	go performWalletActions(client, "user1", wallet.WALLET_CURRENCY_BTC)
	go performWalletActions(client, "user2", wallet.WALLET_CURRENCY_BTC)

	wg.Wait()
}

func performWalletActions(client wallet.WalletClient, userId string, currency wallet.WALLET_CURRENCY) {
	defer wg.Done()

	// Check balance - should return error as no wallet is created
	if _, err := client.GetBalance(context.Background(), &wallet.GetBalanceRequest{UserId: userId, Currency: currency}); err != nil {
		fmt.Printf("GetBalance returned error for user %s & currency %s: %v\n", userId, currency, err)
	}

	// Create wallet for user - should return a successful response
	if resp, err := client.CreateWallet(context.Background(), &wallet.CreateWalletRequest{UserId: userId, Currency: currency}); err != nil {
		fmt.Printf("CreateWallet returned error for user %s & currency %s: %v\n", userId, currency, err)
	} else {
		fmt.Printf("CreateWallet returned response for user %s & currency %s: %v\n", userId, currency, resp)
	}

	// Deposit to wallet for user - should return a successful response
	if resp, err := client.Deposit(context.Background(), &wallet.DepositRequest{UserId: userId, Currency: currency, Amount: 100}); err != nil {
		fmt.Printf("Deposit returned error for user %s & currency %s: %v\n", userId, currency, err)
	} else {
		fmt.Printf("Deposit returned response for user %s & currency %s: %v\n", userId, currency, resp)
	}

	// Withdraw from wallet for user - should return a successful response
	if resp, err := client.Withdraw(context.Background(), &wallet.WithdrawRequest{UserId: userId, Currency: currency, Amount: 80}); err != nil {
		fmt.Printf("Withdraw returned error for user %s & currency %s: %v\n", userId, currency, err)
	} else {
		fmt.Printf("Withdraw returned response for user %s & currency %s: %v\n", userId, currency, resp)
	}

	// Get balance for wallet for user - should return a successful response
	if resp, err := client.GetBalance(context.Background(), &wallet.GetBalanceRequest{UserId: userId, Currency: currency}); err != nil {
		fmt.Printf("GetBalance returned error for user %s & currency %s: %v\n", userId, currency, err)
	} else {
		fmt.Printf("GetBalance returned response for user %s & currency %s: %v\n", userId, currency, resp)
	}

	// Get txn history for wallet for user - should return a successful response
	if resp, err := client.TxnHistory(context.Background(), &wallet.TxnHistoryRequest{UserId: userId, Currency: currency, Page: 1, Size: 10}); err != nil {
		fmt.Printf("TxnHistory returned error for user %s & currency %s: %v\n", userId, currency, err)
	} else {
		fmt.Printf("TxnHistory returned response for user %s & currency %s: %v\n", userId, currency, resp)
	}
}
