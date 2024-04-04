package db

import (
	"context"
	"testing"

	"github.com/koliadervyanko/url-shortener-backend.git/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomString(5),
		Email:    util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
