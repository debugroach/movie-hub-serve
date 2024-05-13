package db

import (
	"context"
	"testing"

	"github.com/debugroach/video-hub-serve/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	username := util.GenerateRandomString(6)
	password := util.GenerateRandomString(6)
	result, err := queries.CreateUser(context.Background(), CreateUserParams{
		username, password,
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	user, err := queries.GetUser(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
	require.Equal(t, password, user.Password)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	userByUsername, err := queries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.ID, userByUsername.ID)
	require.Equal(t, user.Username, userByUsername.Username)
	require.Equal(t, user.Password, userByUsername.Password)
	require.Equal(t, user.CreatedAt, userByUsername.CreatedAt)
}
func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	users, err := queries.ListUsers(context.Background(), ListUsersParams{
		Limit: 5, Offset: 5})
	require.NoError(t, err)
	require.Len(t, users, 5)
}
func TeseDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := queries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)

	userByname, err := queries.GetUser(context.Background(), user.Username)
	require.Error(t, err)
	require.Nil(t, userByname)
}
