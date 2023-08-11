package service

import (
	balance_app "users-balance-monitoring"
	"users-balance-monitoring/pkg/repository"
)

type TransactionHistoryService struct {
	repo repository.TransactionHistory
}

func NewTransactionHistoryService(repo repository.TransactionHistory) *TransactionHistoryService {
	return &TransactionHistoryService{repo: repo}
}

func (s *TransactionHistoryService) GetAllTransactions(userId int, param string) ([]balance_app.TransactionList, int, error) {
	return s.repo.GetAllTransactions(userId, param)
}
func (s *TransactionHistoryService) PaginationTransactions(userId int, param string, limitOfTransactions, offset int) ([]balance_app.TransactionList, error) {
	return s.repo.PaginationTransactions(userId, param, limitOfTransactions, offset)
}
