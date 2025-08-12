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
	conn     *sql.DB
}

func NewClient(host string, port int, database string, username string, password string) (*Client, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database))
	if err != nil {
		return nil, err
	}
	c := &Client{
		host:     host,
		port:     port,
		database: database,
		username: username,
		password: password,
		conn:     conn,
	}
	return c, nil
}

func (c *Client) QueryRow(ctx context.Context, queryTemplate string, args ...any) *sql.Row {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return c.conn.QueryRow(query)
}

func (c *Client) Query(ctx context.Context, queryTemplate string, args ...any) (*sql.Rows, error) {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return c.conn.Query(query)
}

func (c *Client) Exec(ctx context.Context, queryTemplate string, args ...any) (sql.Result, error) {
	query := fmt.Sprintf(queryTemplate, args...)
	tflog.Info(ctx, "MySQL SQL:", map[string]any{"SQL": query})
	return c.conn.Exec(query)
}
