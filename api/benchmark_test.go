package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func BenchmarkCreateAccount(b *testing.B) {
	user, _ := randomUserBench(b)
	account := randomAccountBench(user.Username)

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		CreateAccount(gomock.Any(), gomock.Any()).
		Times(b.N).
		Return(account, nil)

	server := newTestServer(b, store)

	data := gin.H{
		"owner":    user.Username, // Include owner field
		"currency": account.Currency,
	}
	body, err := json.Marshal(data)
	require.NoError(b, err)

	url := "/accounts"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Create new recorder for each iteration
		recorder := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		require.NoError(b, err)

		addAuthorizationBench(b, request, server.tokenMaker, "bearer", user.Username, time.Minute)
		server.router.ServeHTTP(recorder, request)
		require.Equal(b, http.StatusOK, recorder.Code)
	}
}

func BenchmarkGetAccount(b *testing.B) {
	user, _ := randomUserBench(b)
	account := randomAccountBench(user.Username)

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(b.N).
		Return(account, nil)

	server := newTestServer(b, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(b, err)

		addAuthorizationBench(b, request, server.tokenMaker, "bearer", user.Username, time.Minute)
		server.router.ServeHTTP(recorder, request)
		require.Equal(b, http.StatusOK, recorder.Code)
	}
}

func BenchmarkCreateTransfer(b *testing.B) {
	user1, _ := randomUserBench(b)
	user2, _ := randomUserBench(b)
	account1 := randomAccountBench(user1.Username)
	account2 := randomAccountBench(user2.Username)
	account1.Currency = "USD"
	account2.Currency = "USD"

	amount := int64(10)

	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	buildTransferStubs(store, account1, account2, amount, b.N)

	server := newTestServer(b, store)
	recorder := httptest.NewRecorder()

	data := gin.H{
		"from_account_id": account1.ID,
		"to_account_id":   account2.ID,
		"amount":          amount,
		"currency":        "USD",
	}
	body, err := json.Marshal(data)
	require.NoError(b, err)

	url := "/transfers"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		require.NoError(b, err)

		addAuthorizationBench(b, request, server.tokenMaker, "bearer", user1.Username, time.Minute)
		server.router.ServeHTTP(recorder, request)
		require.Equal(b, http.StatusOK, recorder.Code)
	}
}

func buildTransferStubs(store *mockdb.MockStore, account1, account2 db.Account, amount int64, times int) {
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(times).Return(account1, nil)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(times).Return(account2, nil)
	store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(times).Return(db.TransferTXResult{}, nil)
}

// Helper function for testing
func newTestServer(t testing.TB, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

// Benchmark-specific helper functions
func addAuthorizationBench(
	t testing.TB,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("authorization", authorizationHeader)
}

func randomUserBench(t testing.TB) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func randomAccountBench(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
