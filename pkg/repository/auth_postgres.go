package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	balance_app "users-balance-monitoring"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user balance_app.User) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	var balance = 0.00
	createAccountQuery := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := tx.QueryRow(createAccountQuery, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createBalanceQuery := fmt.Sprintf("INSERT INTO %s (user_id, balance) VALUES ($1, $2)", balanceList)
	_, err = tx.Exec(createBalanceQuery, id, balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}
func (r *AuthRepository) GetUser(username, password string) (balance_app.User, error) {
	var user balance_app.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
