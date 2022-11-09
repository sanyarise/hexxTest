package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sanyarise/hezzl/internal/pb"
)

// UserStore is an interface to store user
type UserStore interface {
	SaveUser(ctx context.Context, user *pb.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) (chan *pb.User, error)
}

type UserPostgresStore struct {
	db *sql.DB
}

type PgUser struct {
	Id        string
	CreatedAt time.Time
	DeletedAt *time.Time
	Name      string
}

// ErrAlreadyExists is returned when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

func NewUserPostgresStore(user string, pass string, host string, port string) (*UserPostgresStore, error) {
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres", user, pass, host, port)
	//"postgres://postgres:example@localhost:5432/postgres"
	db, err := sql.Open("pgx", dns)
	if err != nil {
		log.Printf("error on sql open: %s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("error on db ping: %s", err)
		db.Close()
		return nil, err
	}

	us := &UserPostgresStore{
		db: db,
	}
	return us, nil
}

func (ups *UserPostgresStore) Close() {
	ups.db.Close()
}

// Save saves user to the store
func (ups *UserPostgresStore) SaveUser(ctx context.Context, user *pb.User) error {
	pgu := &PgUser{
		Id:        user.Id,
		CreatedAt: time.Now(),
		Name:      user.Name,
	}
	tx, err := ups.db.Begin()
	if err != nil {
		log.Printf("error on begin transaction: %v", err)
		return fmt.Errorf("error on begin transaction: %w", err)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO users
	(id, created_at, name) values ($1, $2, $3)`,
		pgu.Id,
		pgu.CreatedAt,
		pgu.Name,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("error on insert values into table: %v", err)
		return fmt.Errorf("error on insert values into table: %w", err)
	}
	tx.Commit()
	return nil
}

// Delete remove user from store
func (ups *UserPostgresStore) DeleteUser(ctx context.Context, id string) error {
	tx, err := ups.db.Begin()
	if err != nil {
		log.Printf("error on begin transaction: %v", err)
		return fmt.Errorf("error on begin transaction: %w", err)
	}
	_, err = tx.ExecContext(ctx, `UPDATE users SET deleted_at = $2 WHERE id = $1`,
		id, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		log.Printf("error on delete values from table: %v", err)
		return fmt.Errorf("error on delete values from table: %w", err)
	}
	tx.Commit()
	return nil
}

func (ups *UserPostgresStore) GetAllUsers(ctx context.Context) (chan *pb.User, error) {
	chout := make(chan *pb.User, 1000)

	go func() {
		defer close(chout)
		pguser := &PgUser{}

		rows, err := ups.db.QueryContext(ctx, `
		SELECT * FROM users WHERE deleted_at is null`)
		if err != nil {
			log.Printf("error on get all users: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&pguser.Id,
				&pguser.CreatedAt,
				&pguser.DeletedAt,
				&pguser.Name,
			); err != nil {
				log.Printf("error on rows.Scan(): %v", err)
				return
			}

			chout <- &pb.User{
				Id:   pguser.Id,
				Name: pguser.Name,
			}
		}
	}()
	if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
		log.Println("context is cancelled")
		return chout, errors.New("context is cancelled")
	}
	return chout, nil
}
