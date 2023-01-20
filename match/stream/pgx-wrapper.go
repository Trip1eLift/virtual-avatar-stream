package stream

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Docs: https://pkg.go.dev/github.com/jackc/pgx/v5
// Repo: https://github.com/jackc/pgx
// Multi-statements: https://stackoverflow.com/questions/38998267/how-to-execute-a-sql-file

type PGXW struct {
}

type Conn struct {
	*pgx.Conn
}

func (c *Conn) Close() error {
	return c.Conn.Close(context.Background())
}

func (c *Conn) Exec(sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.Conn.Exec(context.Background(), sql, arguments...)
}

func (c *Conn) Query(sql string, args ...any) (pgx.Rows, error) {
	return c.Conn.Query(context.Background(), sql, args...)
}

func (c *Conn) QueryRow(sql string, args ...any) pgx.Row {
	return c.Conn.QueryRow(context.Background(), sql, args...)
}

func (p *PGXW) Connect(psqlurl string) (Conn, error) {
	conn, err := pgx.Connect(context.Background(), psqlurl)
	Conn := Conn{
		conn,
	}
	return Conn, err
}

var pgxw = PGXW{}
