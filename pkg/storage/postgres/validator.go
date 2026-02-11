package postgres

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Validator interface {
	Pool() *pgxpool.Pool
	ValidateQueries(queries []string) error
}

type Config struct {
	URL      string
	Username string
	Password string
	Name     string
}

func (c *Config) connection() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s://%s:%s@%s", "postgres",
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.URL))
	sb.WriteString(fmt.Sprintf("/%s", c.Name))
	return sb.String()
}

type Database struct {
	p *pgxpool.Pool
}

func NewDatabase(cfg *Config) (*Database, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.connection())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Database{
		p: pool,
	}, nil
}

func (d *Database) ValidateQueries(queries []string) error {
	ctx := context.Background()
	conn, err := d.p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for i, q := range queries {
		err = d.validateQuery(ctx, tx.Conn(), q)
		if err == nil {
			continue
		}

		return fmt.Errorf("validate failed for statement [%d]: %w", i, err)
	}

	return nil
}

func (d *Database) Close() {
	d.p.Close()
}

func (d *Database) Pool() *pgxpool.Pool {
	return d.p
}

func (d *Database) Healthy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return errors.Wrap(d.p.Ping(ctx), "postgres is unhealthy")
}

func (d *Database) validateQuery(ctx context.Context, conn *pgx.Conn, query string) error {
	_, err := conn.Prepare(ctx, "validate_query", query)
	if err != nil {
		return fmt.Errorf("failed to create statement: %w", err)
	}

	if err = conn.Deallocate(ctx, "validate_query"); err != nil {
		return fmt.Errorf("failed to close statement: %w", err)
	}

	return nil
}
