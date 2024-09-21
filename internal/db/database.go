package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db *sql.DB
}

type Storage interface {
	CreateUser(*User) error
	GetAccountByEmail(string) (*User, error)
	CreateSession(string, time.Time) error
	GetSessionById(string) (*Session, error)
	DeleteSession(string) error
}

func NewSqliteStore() (*SqliteStore, error) {

	os.Remove("./database/session.db")

	db, err := sql.Open("sqlite3", "./database/session.db")

	fmt.Println("db created ")

	if err != nil {
		return nil, err
	}
	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) Init() error {
	return s.createUserAndSessionTable()
}

func (s *SqliteStore) createUserAndSessionTable() error {

	sqlStmt := `
  create table if not exists users (
  user_id integer primary key,
  username varchar(100),
  password varchar(60)
  );
  `
	_, err := s.db.Exec(sqlStmt)

	if err != nil {
		return err
	}

	sqlStmt = `
  create table if not exists session (
  session_id varchar(40) primary key,
  valid_time text
  );
  `
	_, err = s.db.Exec(sqlStmt)

	return err
}

func (s *SqliteStore) CreateUser(u *User) error {

	_, err := s.db.Exec(`insert into users(username,password) values (?,?);`, u.Username, u.Password)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStore) GetAccountByEmail(username string) (*User, error) {
	var user User
	err := s.db.QueryRow(`select username,password from users where username = $1`, username).Scan(&user.Username, &user.Password)

	if err != nil {
		return nil, err
	}
	return &User{Username: user.Username, Password: user.Password}, nil
}

func (s *SqliteStore) CreateSession(sessionId string, ttx time.Time) error {

	_, err := s.db.Exec(`insert into session(session_id,valid_time) values(?,?);`, sessionId, ttx)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStore) GetSessionById(sessionId string) (*Session, error) {
	var session Session
	err := s.db.QueryRow(`select * from session where session_id = $1`, session).Scan(&session.SessionId, &session.Ttx)

	if err != nil {
		return nil, fmt.Errorf("Session not available please login again")
	}

	return &Session{SessionId: session.SessionId, Ttx: session.Ttx}, nil

}

func (s *SqliteStore) DeleteSession(sessionId string) error {

	_, err := s.db.Exec(`delete from session where session_id = $1`, sessionId)

	if err != nil {
		return err
	}

	return nil
}
