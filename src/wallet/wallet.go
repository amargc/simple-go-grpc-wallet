package wallet

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var SupportedCurrencies = [...]string{"USD", "BTC", "ETC"}

type ledgerRecord struct {
	timestamp int64
	amount    int64
}

type wallet struct {
	mux      *sync.RWMutex
	balances sync.Map // [string]int64
	ledger   map[string][]ledgerRecord
}
type Server struct {
	wallets map[WALLET_CURRENCY]*wallet
}

func NewServer() *Server {
	server := Server{wallets: make(map[WALLET_CURRENCY]*wallet)}
	for i := WALLET_CURRENCY_USD; i < WALLET_CURRENCY_ETC; i++ {
		server.wallets[i] = &wallet{ledger: make(map[string][]ledgerRecord), mux: &sync.RWMutex{}}
	}
	return &server
}

func (s *Server) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*CreateWalletResponse, error) {
	wallet, ok := s.wallets[req.Currency]
	if !ok {
		return nil, fmt.Errorf("%s wallet is not available", SupportedCurrencies[req.Currency])
	}
	if _, ok := wallet.balances.Load(req.UserId); ok {
		return nil, fmt.Errorf("user %s already has a %s wallet", req.UserId, SupportedCurrencies[req.Currency])
	}
	var balance int64 = 0
	wallet.balances.Store(req.UserId, balance)
	return &CreateWalletResponse{UserId: req.UserId, Balance: balance, Currency: req.Currency}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *GetBalanceRequest) (*GetBalanceResponse, error) {
	wallet, ok := s.wallets[req.Currency]
	if !ok {
		return nil, fmt.Errorf("%s wallet is not available", SupportedCurrencies[req.Currency])
	}
	balance_i, ok := wallet.balances.Load(req.UserId)
	if !ok {
		return nil, fmt.Errorf("user %s does not have a %s wallet", req.UserId, SupportedCurrencies[req.Currency])
	}
	balance := balance_i.(int64)
	return &GetBalanceResponse{Balance: balance, Currency: req.Currency}, nil
}

func (s *Server) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	wallet, ok := s.wallets[req.Currency]
	if !ok {
		return nil, fmt.Errorf("%s wallet is not available", SupportedCurrencies[req.Currency])
	}
	var balance int64
	balance_i, ok := wallet.balances.Load(req.UserId)
	if ok {
		balance = balance_i.(int64)
	}
	balance += req.Amount
	wallet.balances.Store(req.UserId, balance)
	wallet.mux.Lock()
	wallet.ledger[req.UserId] = append(wallet.ledger[req.UserId], ledgerRecord{timestamp: time.Now().Unix(), amount: req.Amount})
	wallet.mux.Unlock()
	return &DepositResponse{Balance: balance, Currency: req.Currency}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *WithdrawRequest) (*WithdrawResponse, error) {
	wallet, ok := s.wallets[req.Currency]
	if !ok {
		return nil, fmt.Errorf("%s wallet is not available", SupportedCurrencies[req.Currency])
	}
	var balance int64
	balance_i, ok := wallet.balances.Load(req.UserId)
	if !ok {
		return nil, fmt.Errorf("user %s does not have a %s wallet", req.UserId, SupportedCurrencies[req.Currency])
	}
	balance = balance_i.(int64)
	if balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}
	balance -= req.Amount
	wallet.balances.Store(req.UserId, balance)
	wallet.mux.Lock()
	wallet.ledger[req.UserId] = append(wallet.ledger[req.UserId], ledgerRecord{timestamp: time.Now().Unix(), amount: -req.Amount})
	wallet.mux.Unlock()
	return &WithdrawResponse{Balance: balance, Currency: req.Currency}, nil
}

func (s *Server) TxnHistory(ctx context.Context, req *TxnHistoryRequest) (*TxnHistoryResponse, error) {
	var result *TxnHistoryResponse
	if req.Page < 1 {
		return nil, fmt.Errorf("page %d is not valid input", req.Page)
	}
	if req.Size < 1 {
		return nil, fmt.Errorf("size %d is not valid input", req.Size)
	}
	minIndex, maxIndex := req.Size*(req.Page-1), req.Size*req.Page-1

	wallet, ok := s.wallets[req.Currency]
	if !ok {
		return nil, fmt.Errorf("%s wallet is not available", SupportedCurrencies[req.Currency])
	}
	if _, ok := wallet.balances.Load(req.UserId); !ok {
		return nil, fmt.Errorf("user %s does not have a %s wallet", req.UserId, SupportedCurrencies[req.Currency])
	}

	wallet.mux.RLock()
	ledger, ok := wallet.ledger[req.UserId]
	if !ok {
		return nil, fmt.Errorf("could not get txns for user %s and %s wallet", req.UserId, SupportedCurrencies[req.Currency])
	}
	lenLedger := len(ledger)
	var data []*TxnHistoryResponse_TxnRecord
	if lenLedger > int(maxIndex) {
		for _, record := range ledger[minIndex:maxIndex] {
			data = append(data, &TxnHistoryResponse_TxnRecord{Timestamp: record.timestamp, Amount: record.amount})
		}
	} else if lenLedger > int(minIndex) {
		maxIndex = int32(lenLedger)
		for _, record := range ledger[minIndex:maxIndex] {
			data = append(data, &TxnHistoryResponse_TxnRecord{Timestamp: record.timestamp, Amount: record.amount})
		}
	} else {
		return nil, fmt.Errorf("no more data is available")
	}
	defer wallet.mux.RUnlock()

	result = &TxnHistoryResponse{UserId: req.UserId, Currency: req.Currency, Total: int32(lenLedger), Page: req.Page, Size: req.Size, Data: data}
	return result, nil
}
