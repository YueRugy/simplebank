package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	pwd := RandomOwner()
	hashPassword, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)
	require.NoError(t, CheckPassword(pwd, hashPassword))
	wrongPwd := RandomOwner()
	require.EqualError(t, CheckPassword(wrongPwd, hashPassword),
		bcrypt.ErrMismatchedHashAndPassword.Error())
	hashPassword2, err := HashPassword(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword, hashPassword2)
}
