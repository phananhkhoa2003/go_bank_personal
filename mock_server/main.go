package main

import (
	"fmt"
	"log"
	"simple_bank/api"
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"time"

	"github.com/golang/mock/gomock"
)

func main() {
	// Load config
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Create mock controller
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	// Create mock store
	mockStore := mockdb.NewMockStore(ctrl)

	// Setup common mock expectations that will be used for load testing
	setupMockExpectations(mockStore)

	// Create server with mock store
	server, err := api.NewServer(config, mockStore)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	fmt.Println("Starting mock server for load testing on :8081...")
	fmt.Println("Mock database configured with sample data")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST /users - Create user")
	fmt.Println("  POST /users/login - Login user")
	fmt.Println("  GET /accounts - List accounts (auth required)")
	fmt.Println("  POST /accounts - Create account (auth required)")
	fmt.Println("  POST /transfers - Create transfer (auth required)")
	fmt.Println("")
	fmt.Println("To run load test against this mock server:")
	fmt.Println("  1. Update loadtest/main.go BaseURL to http://localhost:8081")
	fmt.Println("  2. Run: make loadtest")

	err = server.Start(":8081")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func setupMockExpectations(store *mockdb.MockStore) {
	// Create a proper hashed password for "secret123"
	hashedPassword, _ := util.HashPassword("secret123")

	// Mock user creation - always succeeds
	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(db.User{
			Username:  "testuser",
			FullName:  "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
		}, nil)

	// Mock get user by username - always returns a valid user with proper hash
	store.EXPECT().
		GetUser(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(db.User{
			Username:          "testuser",
			HashedPassword:    hashedPassword,
			FullName:          "Test User",
			Email:             "test@example.com",
			PasswordChangedAt: time.Now().Add(-time.Hour),
			CreatedAt:         time.Now().Add(-time.Hour),
		}, nil)

	// Mock account creation - always succeeds
	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(db.Account{
			ID:       1,
			Owner:    "testuser",
			Balance:  100,
			Currency: "USD",
		}, nil)

	// Mock get account - always returns valid account
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(db.Account{
			ID:       1,
			Owner:    "testuser",
			Balance:  100,
			Currency: "USD",
		}, nil)

	// Mock list accounts - returns sample accounts
	store.EXPECT().
		ListAccounts(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return([]db.Account{
			{
				ID:       1,
				Owner:    "testuser",
				Balance:  100,
				Currency: "USD",
			},
			{
				ID:       2,
				Owner:    "testuser",
				Balance:  200,
				Currency: "EUR",
			},
		}, nil)

	// Mock transfer transaction - always succeeds
	store.EXPECT().
		TransferTx(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(db.TransferTXResult{
			Transfer: db.Transfer{
				ID:            1,
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        10,
			},
			FromAccount: db.Account{
				ID:       1,
				Owner:    "testuser",
				Balance:  90,
				Currency: "USD",
			},
			ToAccount: db.Account{
				ID:       2,
				Owner:    "testuser",
				Balance:  210,
				Currency: "USD",
			},
			FromEntry: db.Entry{
				ID:        1,
				AccountID: 1,
				Amount:    -10,
			},
			ToEntry: db.Entry{
				ID:        2,
				AccountID: 2,
				Amount:    10,
			},
		}, nil)
}
