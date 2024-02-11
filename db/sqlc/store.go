package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model"
)

// Store provides all functions to execute db queries and transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	EntryTx(ctx context.Context, arg EntryTxParams) (EntryTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

/*	
	* Method: execTx
		- executes a function within a database transaction
 */
func (store *SQLStore) execTx (ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	
	return tx.Commit()
}

/*
	* Method: TransferTx
		- a concurrent db transaction that performs a money transfer from one account to the other
		- it creates a transfer record, creat account entries and updates the accounts' balances
*/
type TransferTxParams struct {
	FromAccountID	int64 `json:"from_account_id"`
	ToAccountID 	int64 `json:"to_account_id"`
	Amount			int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer 		Transfer	`json:"transfer"`	
	FromAccount		Account		`json:"from_account"`
	ToAccount		Account		`json:"to_account"`
	FromEntry		Entry		`json:"from_entry"`
	ToEntry			Entry		`json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error{
		var err error

		// create a transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// create entries for each account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		
		// update accounts' balances
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, 
				q, 
				arg.FromAccountID, 
				-arg.Amount, 
				arg.ToAccountID, 
				arg.Amount,
			)
			} else {
				result.FromAccount, result.ToAccount, err = addMoney(
					ctx, 
					q, 
					arg.ToAccountID, 
					arg.Amount,
					arg.FromAccountID, 
					-arg.Amount, 
				)
				
			}	
		return nil
	})

	return result, err
}

/* 
	* Method: addmoney
		- it serves the TransferTx function by updating the accounts' balances
			for both the sender and the receiver
 */
func addMoney (
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})

	return
}

/* 
	* Method: EntryTx
		todo : a concurrent db transaction that would performs adding money to an account
		todo : it creates an entry and update the account's balance
 */

 type EntryTxParams struct {
	AccountID	int64	`json:"account_id"`
	Amount		int64	`json:"amount"`
 }

 type EntryTxResult struct {
	Entry	Entry	`json:"entry"`
	Account	Account	`json:"account"`
 }

func (store *SQLStore) EntryTx(ctx context.Context, arg EntryTxParams) (EntryTxResult, error) {
	var result EntryTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create entry
		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		// update account's balance
		result.Account, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		
		return nil
	})

	return result, err
 }