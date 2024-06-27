package db

import (
    "database/sql"
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
    query := `UPDATE users_states SET state = $1 WHERE telegram_id = $2`
    _, err := db.Conn.Exec(query, state, tgID)
    return err
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
