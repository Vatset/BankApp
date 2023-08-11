package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	balance_app "users-balance-monitoring"
)

type TransactionHistoryPostgres struct {
	db *sqlx.DB
}

func NewTransactionHistoryPostgres(db *sqlx.DB) *TransactionHistoryPostgres {
	return &TransactionHistoryPostgres{db: db}
}

func (r *TransactionHistoryPostgres) GetAllTransactions(userId int, param string) ([]balance_app.TransactionList, int, error) {
	var transactions []balance_app.TransactionList
	query := fmt.Sprintf("SELECT id, amount, description, date FROM %s WHERE user_id=$1", transactionPatternTable)
	err := r.db.Select(&transactions, query, userId)
	if err != nil {
		return nil, 0, err
	}
	query = fmt.Sprintf("SELECT  COUNT(*)  FROM %s WHERE user_id=$1"+param, transactionPatternTable)
	var countOfTransactions int
	err = r.db.QueryRow(query, userId).Scan(&countOfTransactions)
	if err != nil {
		return nil, 0, err
	}

	return transactions, countOfTransactions, err
}

func (r *TransactionHistoryPostgres) PaginationTransactions(userId int, param string, limitOfTransactions, offset int) ([]balance_app.TransactionList, error) {
	var transactionsWithPagination []balance_app.TransactionList
	query := fmt.Sprintf("SELECT id, amount, description, date FROM %s WHERE user_id = $1"+param+" LIMIT %d OFFSET %d", transactionPatternTable, limitOfTransactions, offset)
	err := r.db.Select(&transactionsWithPagination, query, userId)
	if err != nil {
		return nil, err
	}

	return transactionsWithPagination, nil
}
