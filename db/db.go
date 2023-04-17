package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func Connect(connString string) (*DB, error) {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	return &DB{conn: conn}, nil
}

func (db *DB) Close() {
	db.conn.Close()
}

func (db *DB) CreateUser(userID string) error {
	_, err := db.conn.Exec(`INSERT INTO users (user_id, dpk_balance) VALUES ($1, 0) ON CONFLICT (user_id) DO NOTHING`, userID)
	return err
}

func (db *DB) UpdateUser(userID string, dpkBalance int) error {
	_, err := db.conn.Exec(`UPDATE users SET dpk_balance = $2 WHERE user_id = $1`, userID, dpkBalance)
	return err
}

func (db *DB) GetUser(userID string) (int, error) {
	var dpkBalance int
	err := db.conn.QueryRow(`SELECT dpk_balance FROM users WHERE user_id = $1`, userID).Scan(&dpkBalance)
	if err != nil {
		return 0, fmt.Errorf("error retrieving user: %v", err)
	}
	return dpkBalance, nil
}

func (db *DB) CreateEvent(eventID, name, description string, startTime, endTime int64, dkpReward int) error {
	_, err := db.conn.Exec(`INSERT INTO events (event_id, name, description, start_time, end_time, dkp_reward) VALUES ($1, $2, $3, $4, $5, $6)`, eventID, name, description, startTime, endTime, dkpReward)
	return err
}

func (db *DB) UpdateEvent(eventID, name, description string, startTime, endTime int64, dkpReward int) error {
	_, err := db.conn.Exec(`UPDATE events SET name = $2, description = $3, start_time = $4, end_time = $5, dkp_reward = $6 WHERE event_id = $1`, eventID, name, description, startTime, endTime, dkpReward)
	return err
}

func (db *DB) GetEvent(eventID string) (string, string, int64, int64, int, error) {
	var name, description string
	var startTime, endTime int64
	var dkpReward int
	err := db.conn.QueryRow(`SELECT name, description, start_time, end_time, dkp_reward FROM events WHERE event_id = $1`, eventID).Scan(&name, &description, &startTime, &endTime, &dkpReward)
	if err != nil {
		return "", "", 0, 0, 0, fmt.Errorf("error retrieving event: %v", err)
	}
	return name, description, startTime, endTime, dkpReward, nil
}

func (db *DB) AddUserAttendence(userID, eventID string) error {
	_, err := db.conn.Exec(`INSERT INTO user_attendence (user_id, event_id) VALUES ($1, $2) ON CONFLICT (user_id, event_id) DO NOTHING`, userID, eventID)
	return err
}
