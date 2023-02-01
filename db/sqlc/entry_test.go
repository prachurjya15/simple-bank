package db

import (
	"context"
	"testing"

	"github.com/prachurjya15/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func generateRandomEntry(t *testing.T, a *Account) Entry {
	arg := CreateEntryParams{
		a.ID,
		util.RandomBalance(),
	}
	e1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, e1)
	require.Equal(t, e1.Amount, arg.Amount)
	require.Equal(t, e1.AccountID, arg.AccountID)

	require.NotZero(t, e1.ID)
	require.NotZero(t, e1.CreatedAt)
	return e1
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	generateRandomEntry(t, &account)
}

func TestGetEntryById(t *testing.T) {
	account := createRandomAccount(t)
	e1 := generateRandomEntry(t, &account)
	e2, err := testQueries.GetEntryById(context.Background(), e1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, e1)
	require.Equal(t, e1.Amount, e2.Amount)
	require.Equal(t, e1.AccountID, e2.AccountID)
}
