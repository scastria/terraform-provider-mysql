package client

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (c *Client) DbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.username, c.password, c.host, c.port, c.db))
	return db, err
}

func (c *Client) QueryRow(ctx context.Context, db *sql.DB, queryTemplate string, args ...any) *sql.Row {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return db.QueryRow(query)
}

func (c *Client) Query(ctx context.Context, db *sql.DB, queryTemplate string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return db.Query(query)
}

func (c *Client) Exec(ctx context.Context, db *sql.DB, queryTemplate string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return db.Exec(query)
}
