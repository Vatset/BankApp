package repository

import (
	"github.com/jmoiron/sqlx"
	balance_app "users-balance-monitoring"
)

type Authorization interface {
	CreateUser(user balance_app.User) (int, error)
	GetUser(username, password string) (balance_app.User, error) //if there is a user->generate token
}

type BalanceManipulations interface {
	Deposit(pattern balance_app.TransactionList) (float64, error)
	Withdraw(pattern balance_app.TransactionList) (float64, error)
	GetBalance(userId int) (float64, error)
	Transfer(userId int, pattern balance_app.TransactionList) error
	UserChecker(pattern balance_app.TransactionList) error
}

type TransactionHistory interface {
	GetAllTransactions(userId int, param string) ([]balance_app.TransactionList, int, error)
	PaginationTransactions(userId int, param, order string, limitOfTransactions, offset int) ([]balance_app.TransactionList, error)
}

type Repository struct {
	Authorization
	BalanceManipulations
	TransactionHistory
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:        NewAuthRepository(db),
		BalanceManipulations: NewBalancePostres(db),
		TransactionHistory:   NewTransactionHistoryPostgres(db),
	}
}
