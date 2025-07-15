package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkHashPassword(b *testing.B) {
	password := RandomString(8)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		hashedPassword, err := HashPassword(password)
		require.NoError(b, err)
		require.NotEmpty(b, hashedPassword)
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := RandomString(8)
	hashedPassword, err := HashPassword(password)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := CheckPassword(password, hashedPassword)
		require.NoError(b, err)
	}
}

func BenchmarkCheckPasswordWrong(b *testing.B) {
	password := RandomString(8)
	wrongPassword := RandomString(8)
	hashedPassword, err := HashPassword(password)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := CheckPassword(wrongPassword, hashedPassword)
		require.Error(b, err)
	}
}

func BenchmarkRandomString(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		str := RandomString(10)
		require.NotEmpty(b, str)
		require.Len(b, str, 10)
	}
}

func BenchmarkRandomInt(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		num := RandomInt(1, 1000)
		require.GreaterOrEqual(b, num, int64(1))
		require.LessOrEqual(b, num, int64(1000))
	}
}

func BenchmarkRandomOwner(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		owner := RandomOwner()
		require.NotEmpty(b, owner)
		require.Len(b, owner, 6)
	}
}

func BenchmarkRandomMoney(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		money := RandomMoney()
		require.GreaterOrEqual(b, money, int64(0))
		require.LessOrEqual(b, money, int64(1000))
	}
}

func BenchmarkRandomCurrency(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		currency := RandomCurrency()
		require.NotEmpty(b, currency)
		require.Contains(b, []string{"USD", "EUR", "CAD", "VND", "JPY", "AUD"}, currency)
	}
}

func BenchmarkRandomEmail(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		email := RandomEmail()
		require.NotEmpty(b, email)
		require.Contains(b, email, "@")
		require.Contains(b, email, ".com")
	}
}
