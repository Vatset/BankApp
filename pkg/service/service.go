package service

import (
	balance_app "users-balance-monitoring"
	"users-balance-monitoring/pkg/repository"
)

type Authorization interface {
	CreateUser(user balance_app.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type BalanceManipulations interface {
	GetBalance(userId int) (float64, error)
	Deposit(pattern balance_app.TransactionList) (float64, error)
	Withdraw(pattern balance_app.TransactionList) (float64, error)
	Transfer(userId int, pattern balance_app.TransactionList) error
	UserChecker(pattern balance_app.TransactionList) error
}

type TransactionHistory interface {
	PaginationTransactions(userId int, param, order string, limitOfTransactions, offset int) ([]balance_app.TransactionList, error)
	GetAllTransactions(userId int, param string) ([]balance_app.TransactionList, int, error)
}

type Service struct {
	Authorization
	BalanceManipulations
	TransactionHistory
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:        NewAuthService(repos.Authorization),
		BalanceManipulations: NewBalanceService(repos.BalanceManipulations),
		TransactionHistory:   NewTransactionHistoryService(repos.TransactionHistory),
	}
}
