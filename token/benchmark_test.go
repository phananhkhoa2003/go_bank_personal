package token

import (
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BenchmarkJWTCreateToken(b *testing.B) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(b, err)

	username := util.RandomOwner()
	duration := time.Minute

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		token, err := maker.CreateToken(username, duration)
		require.NoError(b, err)
		require.NotEmpty(b, token)
	}
}

func BenchmarkJWTVerifyToken(b *testing.B) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(b, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		payload, err := maker.VerifyToken(token)
		require.NoError(b, err)
		require.NotEmpty(b, payload)
	}
}

func BenchmarkPASETOCreateToken(b *testing.B) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(b, err)

	username := util.RandomOwner()
	duration := time.Minute

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		token, err := maker.CreateToken(username, duration)
		require.NoError(b, err)
		require.NotEmpty(b, token)
	}
}

func BenchmarkPASETOVerifyToken(b *testing.B) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(b, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		payload, err := maker.VerifyToken(token)
		require.NoError(b, err)
		require.NotEmpty(b, payload)
	}
}

func BenchmarkJWTVsPASETO(b *testing.B) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(b, err)

	pasetoMaker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(b, err)

	username := util.RandomOwner()
	duration := time.Minute

	b.Run("JWT_Create", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			token, err := jwtMaker.CreateToken(username, duration)
			require.NoError(b, err)
			require.NotEmpty(b, token)
		}
	})

	b.Run("PASETO_Create", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			token, err := pasetoMaker.CreateToken(username, duration)
			require.NoError(b, err)
			require.NotEmpty(b, token)
		}
	})

	// Create tokens for verification benchmark
	jwtToken, err := jwtMaker.CreateToken(username, duration)
	require.NoError(b, err)

	pasetoToken, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(b, err)

	b.Run("JWT_Verify", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			payload, err := jwtMaker.VerifyToken(jwtToken)
			require.NoError(b, err)
			require.NotEmpty(b, payload)
		}
	})

	b.Run("PASETO_Verify", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			payload, err := pasetoMaker.VerifyToken(pasetoToken)
			require.NoError(b, err)
			require.NotEmpty(b, payload)
		}
	})
}
