package db

import (
	"context"
	"testing"

	"github.com/koliadervyanko/url-shortener-backend.git/util"
	"github.com/stretchr/testify/require"
)

func createRandomUrl(t *testing.T, user User) Url {

	arg := CreateUrlParams{
		Url:   util.RandomUrl(),
		Code:  util.RandomString(5),
		Owner: &user.Username,
	}
	url, err := testStore.CreateUrl(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, url)

	require.Equal(t, arg.Url, url.Url)
	require.Equal(t, arg.Code, url.Code)
	require.Equal(t, arg.Owner, url.Owner)
	return url
}

func TestCreateUrl(t *testing.T) {
	user := createRandomUser(t)
	createRandomUrl(t, user)
}

func TestGetUrlByCode(t *testing.T) {
	user := createRandomUser(t)
	url1 := createRandomUrl(t, user)
	url2, err := testStore.GetUrlByCode(context.Background(), url1.Code)
	require.NoError(t, err)
	require.NotEmpty(t, url2)

	require.Equal(t, url1.Code, url2.Code)
	require.Equal(t, url1.Owner, url2.Owner)
	require.Equal(t, url1.Url, url2.Url)
}

func TestListUrlsByUser(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 5; i++ {
		createRandomUrl(t, user)
	}

	urls, err := testStore.ListUrlsByUser(context.Background(), &user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, urls)

	for _, url := range urls {
		require.Equal(t, *url.Owner, user.Username)
	}
}

func TestUpdateClicks(t *testing.T) {
	user := createRandomUser(t)
	url1 := createRandomUrl(t, user)

	url2, err := testStore.UpdateClicks(context.Background(), url1.Code)
	require.NoError(t, err)
	require.NotEmpty(t, url2)

	require.Equal(t, url1.Code, url2.Code)
	require.Equal(t, url1.Owner, url2.Owner)
	require.Equal(t, url1.Url, url2.Url)
	require.Equal(t, url1.Clicks, url2.Clicks-1)
}
