// Package postgres implements postgres connection.
package postgres

import (
	"context"
	"io/ioutil"

	"github.com/go-pg/pg/v10"
)

// Postgres -.
type Postgres struct {
	DB *pg.DB
}

// New -.
func New(url string) (*Postgres, error) {
	// To connect to a database
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)
	// To check if database is up and running
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}

func (pg *Postgres) Migrate(filePath string) error {
	c, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	sql := string(c)
	if _, err := pg.DB.Exec(sql); err != nil {
		return err
	}
	return nil
}