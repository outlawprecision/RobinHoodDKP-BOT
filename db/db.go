package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

type Event struct {
	EventID     string
	Name        string
	Description string
	StartTime   string
	EndTime     string
	DKP_Reward  int
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) UserExists(userID int64) (bool, error) {
	sqlStatement := `
	SELECT EXISTS (
		SELECT 1
		FROM users
		WHERE user_id = $1
	);`

	var exists bool
	err := db.QueryRow(sqlStatement, userID).Scan(&exists)
	return exists, err
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

func (db *DB) CreateEvent(event *Event) error {
	_, err := db.Exec(`
		INSERT INTO Events (event_id, name, description, start_time, end_time, dkp_reward)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (event_id) DO UPDATE
		SET name = $2, description = $3, start_time = $4, end_time = $5, dkp_reward = $6;
	`, event.EventID, event.Name, event.Description, event.StartTime, event.EndTime, 0) // Initially set dkp_reward to 0
	return err
}

func (db *DB) GetEvent(eventID string) (Event, error) {
	sqlStatement := `
	SELECT event_id, name, description, start_time, end_time, dkp_reward
	FROM events
	WHERE event_id = $1;`

	var event Event
	err := db.QueryRow(sqlStatement, eventID).Scan(&event.EventID, &event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.DKP_Reward)
	return event, err
}

func (db *DB) UpdateEvent(event Event) error {
	sqlStatement := `
	UPDATE events
	SET name = $2, description = $3, start_time = $4, end_time = $5, dkp_reward = $6
	WHERE event_id = $1;`

	_, err := db.Exec(sqlStatement, event.EventID, event.Name, event.Description, event.StartTime, event.EndTime, event.DKP_Reward)
	return err
}

func (db *DB) DeleteEvent(eventID string) error {
	sqlStatement := `
	DELETE FROM events
	WHERE event_id = $1;`

	_, err := db.Exec(sqlStatement, eventID)
	return err
}
