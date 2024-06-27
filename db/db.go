package db

import (
    "database/sql"
    "errors"
    _ "github.com/lib/pq"
    "log"
)

type Database struct {
    Conn *sql.DB
}

func ConnectDB(connectingString string) (*Database, error) {
	conn, err := sql.Open("postgres", connectingString)
    onFail("Failed to open dbconnection %v", err)
    
    err = conn.Ping()
    if err != nil {
        return nil, err
    }
    
    return &Database{Conn: conn}, nil
}

func (db *Database) SetUserState(tgID int64, state string) error  {
    query := `
        INSERT INTO users_states (telegram_id, state)
        VALUES ($1, $2)
        ON CONFLICT (telegram_id)
        DO UPDATE SET state = EXCLUDED.state`
    _, err := db.Conn.Exec(query, tgID, state)
    return err
}

type User struct {
    ID int64
    TelegramID int64
    FullName string
    Phone string
}

func (db *Database) GetUserByTgId(tgID int64) (*User, error) {
    var user User
    query := `SELECT id, telegram_id, fullname, phone FROM users WHERE telegram_id = $1`
    err := db.Conn.QueryRow(query, tgID).Scan(&user.ID, &user.TelegramID, &user.FullName, &user.Phone)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (db *Database) AddUser(tgID int64, fullname string, phone string) error {
    query := `INSERT INTO users (telegram_id, fullname, phone) VALUES ($1,$2,$3)`
    _, err := db.Conn.Exec(query, tgID, fullname, phone)
    return err
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
