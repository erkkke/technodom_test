package postgres

import (
	"context"
	"errors"
	"github.com/erkkke/technodom_test/db"
	"github.com/jmoiron/sqlx"
	"time"
)

const connTimeout = 10

type database struct {
	connURL string
	driver  string

	conn *sqlx.DB
}

func (d *database) Connect() error {
	if d.connURL == "" || d.driver == "" {
		return errors.New("cannot connect to db")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connTimeout)*time.Second)
	defer cancel()

	conn, err := sqlx.Open(d.driver, d.connURL)
	if err != nil {
		return err
	}

	if err = conn.PingContext(ctx); err != nil {
		return err
	}

	d.conn = conn

	return nil
}

func (d *database) Close() error {
	return d.conn.Close()
}

func New(connURL, driver string) db.Store {
	return &database{
		connURL: connURL,
		driver:  driver,
	}
}

func (d *database) RedirectList() ([]db.Redirect, error) {
	//TODO implement me
	panic("implement me")
}

func (d *database) Redirect() (*db.Redirect, error) {
	//TODO implement me
	panic("implement me")
}

func (d *database) CreateRedirect(redirect db.Redirect) (*db.Redirect, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
