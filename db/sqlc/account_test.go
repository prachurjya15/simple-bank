package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/prachurjya15/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	testUser := CreateUserParams{
		util.RandomOwner(),
		"secret",
		util.RandomOwner(),
		util.RandomEmail(),
	}
	createdUser, err := testQueries.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	testAcc := CreateAccountParams{
		createdUser.Username,
		util.RandomBalance(),
		util.RandomCurrency(),
	}
	createdAcc, err := testQueries.CreateAccount(context.Background(), testAcc)
	require.NoError(t, err)
	require.Equal(t, testAcc.Owner, createdAcc.Owner)
	require.Equal(t, testAcc.Balance, createdAcc.Balance)
	require.Equal(t, testAcc.Currency, createdAcc.Currency)
	require.NotNil(t, createdAcc.ID)
	require.NotNil(t, createdAcc.CreatedAt)
	return createdAcc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccountById(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccountById(context.Background(), acc1.ID)
	require.NotEmpty(t, acc2)

	require.NoError(t, err)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	arg := UpdateBalanceParams{
		acc1.ID,
		util.RandomBalance(),
	}
	acc2, err := testQueries.UpdateBalance(context.Background(), arg)
	require.NotEmpty(t, acc2)

	require.NoError(t, err)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccountById(context.Background(), acc1.ID)

	require.Empty(t, acc2)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

}

func TestGetAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	arg := GetAccountsParams{
		lastAccount.Owner,
		5,
		0,
	}
	accounts, err := testQueries.GetAccounts(context.Background(), arg)

	require.NotEmpty(t, accounts)
	require.NoError(t, err)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}
