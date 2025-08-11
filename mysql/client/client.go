package client

import (
	"context"
	"database/sql"
	"fmt"
)

type Client struct {
	host     string
	port     int
	db       string
	username string
	password string
}

func NewClient(host string, port int, db string, username string, password string) (*Client, error) {
	c := &Client{
		host:     host,
		port:     port,
		db:       db,
		username: username,
		password: password,
	}
	return c, nil
}

func (c *Client) DbConnection(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.username, c.password, c.host, c.port, c.db))
	return db, err
}
