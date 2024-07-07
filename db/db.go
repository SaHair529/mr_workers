package db

import (
	"database/sql"
	"errors"
	"fmt"
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

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< USER STATE <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type UserState struct {
	TelegramID int64
	State      string
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

func (db *Database) SetUserState(tgID int64, state string) error {
	query := `
        INSERT INTO users_states (telegram_id, state)
        VALUES ($1, $2)
        ON CONFLICT (telegram_id)
        DO UPDATE SET state = EXCLUDED.state`
	_, err := db.Conn.Exec(query, tgID, state)
	return err
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USER STATE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< REQUEST <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type Request struct {
	ID          int64
	TelegramID  int64
	Specialist  string
	City        string
	Description string
	Free        bool
}

func (db *Database) GetFreeRequest(tgID int64) (Request, error) {
	var request Request
	query := `SELECT id, telegram_id, specialist, city, description FROM requests WHERE telegram_id = $1 AND free = true`
	err := db.Conn.QueryRow(query, tgID).Scan(&request.ID, &request.TelegramID, &request.Specialist, &request.City, &request.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Request{}, nil
		}
		return Request{}, err
	}
	return request, nil
}

func (db *Database) GetRequestById(id int64) (Request, error) {
	var request Request
	query := `SELECT id, telegram_id, specialist, city, description, free FROM requests WHERE id = $1`
	err := db.Conn.QueryRow(query, id).Scan(&request.ID, &request.TelegramID, &request.Specialist, &request.City, &request.Description, &request.Free)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Request{}, nil
		}
		return Request{}, err
	}
	return request, nil
}

func (db *Database) DeleteFreerequest(tgID int64) error {
	query := `DELETE FROM requests WHERE telegram_id = $1 AND free = true`
	_, err := db.Conn.Exec(query, tgID)
	return err
}

func (db *Database) CreateFreeRequest(tgID int64) error {
	query := `INSERT INTO requests (telegram_id) VALUES ($1)`
	_, err := db.Conn.Exec(query, tgID)
	return err
}

func (db *Database) SetUnfreeRequest(tgID int64) error {
	query := `UPDATE requests SET free = false WHERE telegram_id = $1 AND free = true`
	_, err := db.Conn.Exec(query, tgID)
	return err
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> REQUEST >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< USER <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type User struct {
	ID         int64
	TelegramID int64
	FullName   string
	Phone      string
}

func (db *Database) GetUserByTgId(tgID int64) (User, error) {
	var user User
	query := `SELECT id, telegram_id, fullname, phone FROM users WHERE telegram_id = $1`
	err := db.Conn.QueryRow(query, tgID).Scan(&user.ID, &user.TelegramID, &user.FullName, &user.Phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

func (db *Database) AddUser(tgID int64, fullname string, phone string) error {
	query := `INSERT INTO users (telegram_id, fullname, phone) VALUES ($1,$2,$3)`
	_, err := db.Conn.Exec(query, tgID, fullname, phone)
	return err
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USER >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< SPECIALITY <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type Speciality struct {
	ID         int64
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

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SPECIALITY >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< WORKER <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

type Worker struct {
	ID         int64
	TelegramID int64
	FullName   string
	Phone      string
	Speciality string
	City       string
}

func (db *Database) GetWorkerByTgId(tgID int64) (Worker, error) {
	var worker Worker
	query := `SELECT id, fullname, phone, city, speciality FROM workers WHERE telegram_id = $1`
	err := db.Conn.QueryRow(query, tgID).Scan(&worker.ID, &worker.FullName, &worker.Phone, &worker.City, &worker.Speciality)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Worker{}, nil
		}
		return Worker{}, err
	}
	return worker, nil
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

func (db *Database) SetWorkerContactData(tgID int64, fullname, phone string) error {
	query := `
        INSERT INTO workers (telegram_id, fullname, phone)
        VALUES ($1, $2, $3)
        ON CONFLICT (telegram_id)
        DO UPDATE SET fullname = EXCLUDED.fullname, phone = EXCLUDED.phone`
	_, err := db.Conn.Exec(query, tgID, fullname, phone)
	return err
}

func (db *Database) GetFreeWorkersByCityAndSpeciality(city string, speciality string) ([]Worker, error) {
	query := "SELECT id, telegram_id, fullname, phone, speciality, city FROM workers WHERE city = $1 AND speciality = $2"
	rows, err := db.Conn.Query(query, city, speciality)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workers []Worker
	for rows.Next() {
		var worker Worker
		if err := rows.Scan(&worker.ID, &worker.TelegramID, &worker.FullName, &worker.Phone, &worker.Speciality, &worker.City); err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return workers, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> WORKER >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func (db *Database) SetRowField(tgID int64, tableName, fieldName, fieldValue string) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (telegram_id, %s)
        VALUES ($1, $2)
        ON CONFLICT (telegram_id)
        DO UPDATE SET %s = EXCLUDED.%s`, tableName, fieldName, fieldName, fieldName)

	_, err := db.Conn.Exec(query, tgID, fieldValue)
	return err
}

func (db *Database) SetFreeRequestField(tgID int64, fieldName, fieldValue string) error {
	query := fmt.Sprintf(`UPDATE requests SET %s = $2 WHERE free = true AND telegram_id = $1`, fieldName)
	_, err := db.Conn.Exec(query, tgID, fieldValue)
	return err
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
