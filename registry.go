package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Grade int

const (
	trainee Grade = iota + 1
	junior
	middle
	senior
)

var grades = map[Grade]string{
	trainee: "trainee",
	junior:  "junior",
	middle:  "middle",
	senior:  "senior",
}

type User struct {
	Id       int
	Name     string
	Surname  string
	Position Grade
	Project  string
}

type Registry struct {
	p pool
}

// Interface for ease of mocking, exposing only used methods of pgxpool.Pool
type pool interface {
	Ping(context.Context) error
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
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
		name, surname, grades[position], project)
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
