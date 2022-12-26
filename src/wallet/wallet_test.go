package wallet

import (
	"context"
	"net"
	"testing"

	grpc "google.golang.org/grpc"
)

func TestWalletService(t *testing.T) {
	// Create a test server and client.
	s := grpc.NewServer()
	RegisterWalletServer(s, NewServer())

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	go s.Serve(lis)
	defer s.Stop()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := NewWalletClient(conn)

	// Test the CreateWallet method.
	resp0, err := client.CreateWallet(context.Background(), &CreateWalletRequest{UserId: "user1", Currency: WALLET_CURRENCY_USD})
	if err != nil {
		t.Errorf("CreateWallet returned error: %v", err)
	}
	if resp0.Balance != 0 {
		t.Errorf("GetBalance returned incorrect balance: %d", resp0.Balance)
	}

	// Test the GetBalance method.
	resp1, err := client.GetBalance(context.Background(), &GetBalanceRequest{UserId: "user1", Currency: WALLET_CURRENCY_USD})
	if err != nil {
		t.Errorf("GetBalance returned error: %v", err)
	}
	if resp1.Balance != 0 {
		t.Errorf("GetBalance returned incorrect balance: %d", resp1.Balance)
	}

	// Test the Deposit method.
	resp2, err := client.Deposit(context.Background(), &DepositRequest{UserId: "user1", Amount: 100, Currency: WALLET_CURRENCY_USD})
	if err != nil {
		t.Errorf("Deposit returned error: %v", err)
	}
	if resp2.Balance != 100 {
		t.Errorf("Deposit returned incorrect balance: %d", resp2.Balance)
	}

	// Test the Withdraw method.
	resp3, err := client.Withdraw(context.Background(), &WithdrawRequest{UserId: "user1", Amount: 50, Currency: WALLET_CURRENCY_USD})
	if err != nil {
		t.Errorf("Withdraw returned error: %v", err)
	}
	if resp3.Balance != 50 {
		t.Errorf("Withdraw returned incorrect balance: %d", resp3.Balance)
	}

	// Test the TxnHistory method.
	resp4, error := client.TxnHistory(context.Background(), &TxnHistoryRequest{UserId: "user1", Currency: WALLET_CURRENCY_USD, Page: 1, Size: 10})
	if error != nil {
		t.Errorf("TxnHistory returned error: %v", err)
	}
	if resp4.Total != 2 {
		t.Errorf("TxnHistory returned incorrect total count: %d", resp4.Total)
	}
	if len(resp4.Data) != 2 {
		t.Errorf("TxnHistory returned incorrect no. of txns: %d", resp4.Total)
	}
}
