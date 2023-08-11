package promo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slices"
	"log"
	"strings"
)

type Registry struct {
	p pool
}

// Interface for ease of mocking, exposing only used methods of pgxpool.Pool
type pool interface {
	Ping(context.Context) error
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func initDB(connString string) (p pool, err error) {
	if p, err = pgxpool.New(context.Background(), connString); err != nil {
		return nil, fmt.Errorf("failed to create pool connection to database %s: %w", connString, err)
	}
	return p, nil
}

func NewRegistry(connString string) (*Registry, error) {
	p, err := initDB(connString)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Registry{p}, nil
}

func (r *Registry) AddUser(name string, surname string, position Grade, project string) (err error) {
	_, err = r.p.Exec(context.Background(),
		"INSERT INTO usr (name, surname, position, project) VALUES ($1, $2, $3, $4)",
		name, surname, dGrades[position], project)
	if err != nil {
		return fmt.Errorf("unable to INSERT INTO usr: %w", err)
	}
	return nil
}

func (r *Registry) DeleteUser(id int) (err error) {
	rp, err := r.p.Exec(context.Background(),
		"DELETE FROM usr WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("unable to DELETE FROM usr: %w", err)
	}
	if rp.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected when attempt to DELETE FROM usr with id: %v", id)
	}
	return nil
}

func (r *Registry) UpdateUser(id int, m map[string]string) (err error) {
	fields := []string{"name", "surname", "position", "project"}
	s := []string{}
	for k, v := range m {
		if !slices.Contains(fields, k) {
			return fmt.Errorf("illegal key in the map")
		}
		s = append(s, fmt.Sprintf("%s='%s'", k, v))
	}
	req := fmt.Sprintf("UPDATE usr SET %s WHERE id=$1", strings.Join(s, ","))
	rp, err := r.p.Exec(context.Background(), req, id)
	if err != nil {
		return fmt.Errorf("unable to UPDATE usr: %w", err)
	}
	if rp.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected while attempting to UPDATE usr with id: %v", id)
	}
	return nil
}

func (r *Registry) GetUser(id int) (*User, error) {
	row := r.p.QueryRow(context.Background(), "SELECT name, surname, position, project FROM usr WHERE id=$1", id)
	u := &User{}
	var pos string
	err := row.Scan(&u.Name, &u.Surname, &pos, &u.Project)
	if err != nil {
		return nil, fmt.Errorf("unable to get user with id %d: %w", id, err)
	}
	u.Position = bGrades[pos]
	u.Id = id
	return u, nil
}

func (r *Registry) GetAllUsers() (*[]User, error) {
	rows, err := r.p.Query(context.Background(), "SELECT * FROM usr")
	if err != nil {
		return nil, fmt.Errorf("unable to SELECT all users FROM usr: %w", err)
	}
	log.Println("Printing rows")
	u := &User{}
	var us []User
	var pos string

	tr, err := pgx.ForEachRow(rows, []any{&u.Id, &u.Name, &u.Surname, &pos, &u.Project}, func() error {
		u.Position = bGrades[pos]
		us = append(us, *u)
		return nil
	})
	if err != nil || int(tr.RowsAffected()) != len(us) {
		return nil, fmt.Errorf("unable to convert request into names list: %w", err)
	}
	return &us, nil
}
