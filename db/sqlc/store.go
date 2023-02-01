package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	Query *Queries
	Db    *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Db:    db,
		Query: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			return fmt.Errorf("error in executing fn %s . Rollback Error : %s", err, rErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountId int64 `json: from_account_id`
	ToAccountId   int64 `json: to_account_id`
	Amount        int64 `json: "amount"`
}

type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

var TxKey = struct{}{}

// Create Transfer Record
// Create 2 Entry one for from and one for to
// Update the account balance of both
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	var err error
	err = store.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(TxKey)

		fmt.Println(txName, "Create a Transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{arg.FromAccountId, arg.ToAccountId, arg.Amount})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create a Entry In From Account")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{arg.FromAccountId, -arg.Amount})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create a Entry In ToAccount")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{arg.ToAccountId, arg.Amount})
		if err != nil {
			return err
		}
		// TODO: Update the account

		if arg.FromAccountId < arg.ToAccountId {
			fmt.Println(txName, "Updat FromAccount")
			result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{ID: arg.FromAccountId, Amount: -arg.Amount})
			if err != nil {
				return err
			}

			fmt.Println(txName, "Update toAccount")
			result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{ID: arg.ToAccountId, Amount: arg.Amount})
			if err != nil {
				return err
			}
		} else {
			fmt.Println(txName, "Update toAccount")
			result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{ID: arg.ToAccountId, Amount: arg.Amount})
			if err != nil {
				return err
			}
			fmt.Println(txName, "Updat FromAccount")
			result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{ID: arg.FromAccountId, Amount: -arg.Amount})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, nil
}
