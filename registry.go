package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connexion struct {
	p pool
}

// Interface for ease of mocking, exposing only used methods of pgxpool.Pool
type pool interface {
	Ping(context.Context) error
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type User struct {
	Id       int
	Name     string
	Surname  string
	Position Grade
	Project  string
}

func initDB(connString string) (p pool, err error) {
	if p, err = pgxpool.New(context.Background(), connString); err != nil {
		return nil, fmt.Errorf("failed to create pool connection to database %s: %w", connString, err)
	}
	return p, nil
}

func NewConnexion(connString string) (*Connexion, error) {
	p, err := initDB(connString)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Connexion{p}, nil
}

func (c *Connexion) AddUser(name string, surname string, position Grade, project string) (err error) {
	cp, err := c.p.Exec(context.Background(),
		"INSERT INTO usr (name, surname, position, project) VALUES ($1, $2, $3, $4)",
		name, surname, grades[position], project)
	if err != nil {
		return fmt.Errorf("unable to INSERT INTO usr: %w", err)
	}
	return nil
}
