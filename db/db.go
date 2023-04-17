package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

func NewDB(connectionString string) (*Database, error) {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &Database{conn: conn}, nil
}

func (db *Database) Close() {
	db.conn.Close()
}

// CreateUser creates a new user in the users table with the given user_id
func (db *Database) CreateUser(user_id int64) error {
	_, err := db.conn.Exec("INSERT INTO users (user_id, dkp_balance) VALUES ($1, 0) ON CONFLICT DO NOTHING", user_id)
	return err
}

// CheckBalance returns the DP balance of the user with the given user_id
func (db *Database) CheckBalance(user_id int64) (int, error) {
	var balance int
	err := db.conn.QueryRow("SELECT dpk_balance FROM users WHERE user_id = $1", user_id).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("error checking balance: %w", err)
	}
	return balance, nil
}

// AddPoints adds the specified amount of points to the user with the given user_id
func (db *Database) AddPoints(user_id int64, points int) error {
	_, err := db.conn.Exec("UPDATE users SET dpk_balance = dpk_balance + $1 WHERE user_id = $2", points, user_id)
	return err
}

// RemovePoints removes the specified amount of points from the user with the given user_id
func (db *Database) RemovePoints(user_id int64, points int) error {
	_, err := db.conn.Exec("UPDATE users SET dpk_balance = dpk_balance - $1 WHERE user_id = $2", points, user_id)
	return err
}
