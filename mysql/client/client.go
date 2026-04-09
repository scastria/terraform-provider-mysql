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
	database string
	username string
	password string
	Conn     *sql.DB
}

func NewClient(host string, port int, database string, username string, password string, maxOpenConnections int, maxIdleConnections int) (*Client, error) {
	Conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database))
	if err != nil {
		return nil, err
	}
	Conn.SetMaxOpenConns(maxOpenConnections)
	Conn.SetMaxIdleConns(maxIdleConnections)
	c := &Client{
		host:     host,
		port:     port,
		database: database,
		username: username,
		password: password,
		Conn:     Conn,
	}
	return c, nil
}

func (c *Client) QueryRow(ctx context.Context, queryTemplate string, args ...any) (string, *sql.Row) {
	var stats = c.Conn.Stats()
	tflog.Error(ctx, "MySQL Stats:", map[string]any{"InUse": stats.InUse, "Idle": stats.Idle, "Open": stats.OpenConnections})
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return query, c.Conn.QueryRowContext(ctx, query)
}

func (c *Client) Query(ctx context.Context, queryTemplate string, args ...any) (string, *sql.Rows, error) {
	var stats = c.Conn.Stats()
	tflog.Error(ctx, "MySQL Stats:", map[string]any{"InUse": stats.InUse, "Idle": stats.Idle, "Open": stats.OpenConnections})
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	rows, err := c.Conn.QueryContext(ctx, query)
	return query, rows, err
}

func (c *Client) Exec(ctx context.Context, queryTemplate string, args ...any) (string, sql.Result, error) {
	var stats = c.Conn.Stats()
	tflog.Error(ctx, "MySQL Stats:", map[string]any{"InUse": stats.InUse, "Idle": stats.Idle, "Open": stats.OpenConnections})
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	result, err := c.Conn.ExecContext(ctx, query)
	return query, result, err
}
