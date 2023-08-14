package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
	balance_app "users-balance-monitoring"
)

type BalancePostres struct {
	db *sqlx.DB
}

func NewBalancePostres(db *sqlx.DB) *BalancePostres {
	return &BalancePostres{db: db}
}
func (r *BalancePostres) GetBalance(userId int) (float64, error) {
	var balance float64
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceList)
	err := r.db.Get(&balance, query, userId)

	return balance, err
}

func (r *BalancePostres) Deposit(pattern balance_app.TransactionList) (float64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	updateBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE user_id = $2 RETURNING balance", balanceList)
	row := tx.QueryRow(updateBalanceQuery, pattern.Amount, pattern.Id)

	var updatedBalance float64
	err = row.Scan(&updatedBalance)
	if err != nil {
		return 0, err
	}

	createTransactionListQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, date) VALUES ($1, $2, $3,$4)", transactionPatternTable)
	_, err = tx.Exec(createTransactionListQuery, pattern.Id, pattern.Amount, pattern.Description, SetCurrentTime())
	if err != nil {
		return 0, err
	}
	return updatedBalance, tx.Commit()
}

func (r *BalancePostres) Transfer(userId int, pattern balance_app.TransactionList) error {
	var balance float64
	query := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceList)
	if err := r.db.Get(&balance, query, userId); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if balance >= pattern.Amount {
		if userId != pattern.Id {

			updateRecipientBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance +$1 WHERE user_id = $2 ", balanceList)
			_, err = tx.Exec(updateRecipientBalanceQuery, pattern.Amount, pattern.Id)
			if err != nil {
				return err
			}

			updateSenderBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE user_id = $2 ", balanceList)
			_, err = tx.Exec(updateSenderBalanceQuery, pattern.Amount, userId)
			if err != nil {
				return err
			}

			createRecipientTransactionListQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, date) VALUES ($1, $2, $3, $4)", transactionPatternTable)
			_, err = tx.Exec(createRecipientTransactionListQuery, pattern.Id, pattern.Amount, pattern.Description, SetCurrentTime())
			if err != nil {
				return err
			}

			withdrawAmount := 0 - pattern.Amount
			description := "transfer to user " + strconv.Itoa(pattern.Id) + " description: " + pattern.Description
			createSenderTransactionListQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, date) VALUES ($1, $2, $3, $4)", transactionPatternTable)
			_, err = tx.Exec(createSenderTransactionListQuery, userId, withdrawAmount, description, SetCurrentTime())
			if err != nil {
				return err
			}
			return tx.Commit()
		} else {
			return errors.New("You can not transfer money to yourself")
		}

	} else {
		return errors.New("Your balance should be more than the transfer amount")
	}
}
func (r *BalancePostres) UserChecker(pattern balance_app.TransactionList) error {
	checkPatternQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id=$1", balanceList)
	var count int
	row := r.db.QueryRow(checkPatternQuery, pattern.Id)
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		return errors.New("Pattern.Id does not exist in the database")
	}
	return nil
}
func SetCurrentTime() time.Time {
	date := time.Now()
	return date
}

func (r *BalancePostres) Withdraw(pattern balance_app.TransactionList) (float64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	getBalanceQuery := fmt.Sprintf("SELECT balance FROM %s WHERE user_id=$1", balanceList)
	row := tx.QueryRow(getBalanceQuery, pattern.Id)
	var currentBalance float64
	if err := row.Scan(&currentBalance); err != nil {
		return 0, err
	}

	if currentBalance >= pattern.Amount {

		updateBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance-$1 WHERE user_id = $2 RETURNING balance", balanceList)

		row = tx.QueryRow(updateBalanceQuery, pattern.Amount, pattern.Id)
		var updatedBalance float64
		err := row.Scan(&updatedBalance)
		if err != nil {
			return 0, err
		}

		withdrawAmount := 0 - pattern.Amount

		createTransactionListQuery := fmt.Sprintf("INSERT INTO %s (user_id, amount, description, date) VALUES ($1, $2, $3,$4)", transactionPatternTable)
		_, err = tx.Exec(createTransactionListQuery, pattern.Id, withdrawAmount, pattern.Description, SetCurrentTime())
		if err != nil {
			return 0, err
		}
		return updatedBalance, tx.Commit()
	} else {
		return 0, errors.New("Withdraw amount must be less then users balance")
	}
}
