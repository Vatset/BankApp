package service

import (
	balance_app "users-balance-monitoring"
	"users-balance-monitoring/pkg/repository"
)

type BalanceService struct {
	repo repository.BalanceManipulations
}

func NewBalanceService(repo repository.BalanceManipulations) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) Deposit(pattern balance_app.TransactionList) (float64, error) {
	return s.repo.Deposit(pattern)
}

func (s *BalanceService) Withdraw(pattern balance_app.TransactionList) (float64, error) {
	return s.repo.Withdraw(pattern)
}

func (s *BalanceService) GetBalance(userId int) (float64, error) {
	return s.repo.GetBalance(userId)
}
func (s *BalanceService) Transfer(userId int, pattern balance_app.TransactionList) error {
	return s.repo.Transfer(userId, pattern)
}
func (s *BalanceService) UserChecker(pattern balance_app.TransactionList) error {
	return s.repo.UserChecker(pattern)
}
