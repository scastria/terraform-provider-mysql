package client

type Client struct {
	host string
	port int
	db   string
}

func NewClient(host string, port int, db string) (*Client, error) {
	c := &Client{
		host: host,
		port: port,
		db:   db,
	}
	return c, nil
}
