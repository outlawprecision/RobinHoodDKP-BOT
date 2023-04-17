package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) CreateUser(userID int64) error {
	sqlStatement := `
	INSERT INTO users (user_id, dkp_balance)
	VALUES ($1, 0)
	ON CONFLICT (user_id) DO NOTHING;`

	_, err := db.Exec(sqlStatement, userID)
	return err
}

func (db *DB) CheckBalance(userID int64) (int, error) {
	sqlStatement := `
	SELECT dkp_balance
	FROM users
	WHERE user_id = $1;`

	var balance int
	err := db.QueryRow(sqlStatement, userID).Scan(&balance)
	return balance, err
}

func (db *DB) AddPoints(userID int64, points int) error {
	sqlStatement := `
	UPDATE users
	SET dkp_balance = dkp_balance + $2
	WHERE user_id = $1;`

	_, err := db.Exec(sqlStatement, userID, points)
	return err
}

func (db *DB) RemovePoints(userID int64, points int) error {
	sqlStatement := `
	UPDATE users
	SET dkp_balance = dkp_balance - $2
	WHERE user_id = $1
	RETURNING dkp_balance;`

	var newBalance int
	err := db.QueryRow(sqlStatement, userID, points).Scan(&newBalance)
	if err != nil {
		return err
	}

	if newBalance < 0 {
		// Revert the change and return an error if the new balance is negative
		sqlStatement = `
		UPDATE users
		SET dkp_balance = dkp_balance + $2
		WHERE user_id = $1;`
		_, err := db.Exec(sqlStatement, userID, points)
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot remove more points than the current balance")
	}

	return nil
}
