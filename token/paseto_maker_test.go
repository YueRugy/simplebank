package token

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
	"simplebank/util"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	username, duration := util.RandomOwner(), time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}
func TestExpiredPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)
	require.NotEmpty(t, maker)
	username, duration := util.RandomOwner(), -time.Minute
	tokenString, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	payload, err := maker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredAt.Error())
	require.Nil(t, payload)
}
