package db

import (
	"context"
	"testing"

	"github.com/koliadervyanko/url-shortener-backend.git/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(5)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomString(5),
		Password: &hashedPassword,
		Color:    util.RandomColor(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
