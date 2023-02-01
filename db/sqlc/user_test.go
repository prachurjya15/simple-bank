package db

import (
	"context"
	"testing"

	"github.com/prachurjya15/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPwd, err := util.CreateHashedPwd(util.RandomString(6))
	require.NoError(t, err)
	testUser := CreateUserParams{
		util.RandomOwner(),
		hashedPwd,
		util.RandomOwner(),
		util.RandomEmail(),
	}
	createdUser, err := testQueries.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	require.Equal(t, testUser.Username, createdUser.Username)
	require.Equal(t, testUser.HashedPassword, createdUser.HashedPassword)
	require.Equal(t, testUser.Email, createdUser.Email)
	require.Equal(t, testUser.FullName, createdUser.FullName)
	require.NotNil(t, createdUser.PasswordChangedAt)
	require.NotNil(t, createdUser.CreatedAt)
	return createdUser
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NotEmpty(t, user2)

	require.NoError(t, err)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
}
