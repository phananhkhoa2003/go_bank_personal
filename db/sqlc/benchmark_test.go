package db

import (
	"context"
	"simple_bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

// Benchmark database operations
func BenchmarkCreateAccount(b *testing.B) {
	// Pre-create users with different currencies to avoid constraint violations
	currencies := []string{"USD", "EUR", "CAD"}
	users := make([]User, b.N)

	// Setup phase: create users (not timed)
	for i := 0; i < b.N; i++ {
		users[i] = createRandomUserForBenchmark(b)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark only account creation
	for i := 0; i < b.N; i++ {
		arg := CreateAccountParams{
			Owner:    users[i].Username,
			Balance:  util.RandomMoney(),
			Currency: currencies[i%len(currencies)],
		}

		account, err := testStore.CreateAccount(context.Background(), arg)
		require.NoError(b, err)
		require.NotEmpty(b, account)
	}
}

func BenchmarkGetAccount(b *testing.B) {
	// Setup: Create an account first
	account := createRandomAccountForBenchmark(b)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		acc, err := testStore.GetAccount(context.Background(), account.ID)
		require.NoError(b, err)
		require.NotEmpty(b, acc)
	}
}

func BenchmarkListAccounts(b *testing.B) {
	// Setup: Create a user and accounts (one per currency to avoid duplicates)
	user := createRandomUserForBenchmark(b)
	currencies := []string{"USD", "EUR", "CAD"}

	// Create exactly 3 accounts (one per currency)
	for _, currency := range currencies {
		arg := CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomMoney(),
			Currency: currency,
		}
		_, err := testStore.CreateAccount(context.Background(), arg)
		require.NoError(b, err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		arg := ListAccountsParams{
			Owner:  user.Username,
			Limit:  5,
			Offset: 0,
		}

		accounts, err := testStore.ListAccounts(context.Background(), arg)
		require.NoError(b, err)
		require.True(b, len(accounts) > 0, "Should return at least one account")
	}
}

func BenchmarkTransferTx(b *testing.B) {
	// Setup: Create two accounts
	account1 := createRandomAccountForBenchmark(b)
	account2 := createRandomAccountForBenchmark(b)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result, err := testStore.TransferTx(context.Background(), TransferTXParams{
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        10,
		})
		require.NoError(b, err)
		require.NotEmpty(b, result)
	}
}

func BenchmarkTransferTxConcurrent(b *testing.B) {
	account1 := createRandomAccountForBenchmark(b)
	account2 := createRandomAccountForBenchmark(b)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := testStore.TransferTx(context.Background(), TransferTXParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        1,
			})
			require.NoError(b, err)
			require.NotEmpty(b, result)
		}
	})
}

func BenchmarkCreateUser(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		hashedPassword, err := util.HashPassword(util.RandomString(6))
		require.NoError(b, err)

		arg := CreateUserParams{
			Username:       util.RandomOwner(),
			HashedPassword: hashedPassword,
			FullName:       util.RandomOwner(),
			Email:          util.RandomEmail(),
		}

		user, err := testStore.CreateUser(context.Background(), arg)
		require.NoError(b, err)
		require.NotEmpty(b, user)
	}
}

// Helper function for benchmarks
func createRandomAccountForBenchmark(t testing.TB) Account {
	// First create a user
	user := createRandomUserForBenchmark(t)

	// Then create an account with that user as owner
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	return account
}

// Helper function to create a user for benchmarks
func createRandomUserForBenchmark(t testing.TB) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}
