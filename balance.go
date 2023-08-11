package balance_app

import (
	"errors"
	"time"
)

type BalanceList struct {
	Id      int     `json:"id"`
	UserId  int     `json:"userId"`
	Balance float64 `json:"balance"`
}

type TransactionList struct {
	Id          int        `json:"id" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	Description string     `json:"description" binding:"required" swag:"optional"`
	Date        CustomTime `json:"date"  binding:"-" swag:"ignore"`
}

func (t *TransactionList) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	time := time.Time(ct)
	return []byte(`"` + time.Format("02-01-2006 15:04:05") + `"`), nil
}
