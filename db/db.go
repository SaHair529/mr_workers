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

type UserState struct {
    TelegramID int64
    State string
}

func (db *Database) GetUserState(tgID int64) (string, error) {
    var state string
    query := `SELECT state FROM users_states WHERE telegram_id = $1`
    err := db.Conn.QueryRow(query, tgID).Scan(&state)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", nil
        }
        return "", err
    }
    return state, nil
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

type Speciality struct {
    ID int64
    Speciality string
}

func (db *Database) GetSpecialityByTitle(title string) (Speciality, error) {
    var speciality Speciality
    query := `SELECT id, speciality FROM specialities WHERE speciality = $1`
    err := db.Conn.QueryRow(query, title).Scan(&speciality.ID, &speciality.Speciality)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Speciality{}, nil
        }
        return Speciality{}, err
    }
    return speciality, nil
}

func (db *Database) GetAllSpecialities() ([]Speciality, error) {
    query := `SELECT id, speciality FROM specialities`
    rows, err := db.Conn.Query(query)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var specialities []Speciality
    for rows.Next() {
        var spec Speciality
        if err := rows.Scan(&spec.ID, &spec.Speciality); err != nil {
            return nil, err
        }
        specialities = append(specialities, spec)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return specialities, nil
}

func (db *Database) SetWorkerSpeciality(tgID int64, speciality string) error {
    query := `
        INSERT INTO workers (telegram_id, speciality)
        VALUES ($1, $2)
        ON CONFLICT (telegram_id)
        DO UPDATE SET speciality = EXCLUDED.speciality`
    _, err := db.Conn.Exec(query, tgID, speciality)
    return err
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
