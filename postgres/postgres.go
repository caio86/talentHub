package postgres

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	talenthub "github.com/caio86/talentHub"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     uint16
	DBName   string
	SSLMode  string
}

func (c *DBConfig) URL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func (c *DBConfig) Validate() error {
	if c.User == "" {
		return talenthub.Errorf(talenthub.EINVALID, "DB Username required")
	}
	if c.Password == "" {
		return talenthub.Errorf(talenthub.EINVALID, "DB Password required")
	}
	if c.Host == "" {
		return talenthub.Errorf(talenthub.EINVALID, "DB Host required")
	}
	if c.Port == 0 {
		return talenthub.Errorf(talenthub.EINVALID, "DB Port required")
	}
	if c.DBName == "" {
		return talenthub.Errorf(talenthub.EINVALID, "DB Name required")
	}
	if c.SSLMode == "" {
		return talenthub.Errorf(talenthub.EINVALID, "DB SSLMode required")
	}

	return nil
}

type DB struct {
	conn *pgxpool.Pool
	*DBConfig
}

func NewDB(config *DBConfig) *DB {
	db := &DB{
		DBConfig: config,
	}

	return db
}

func (db *DB) Connect() error {
	if err := db.Validate(); err != nil {
		return err
	}

	conn, err := pgxpool.New(context.Background(), db.URL())
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "%s", err)
	}

	db.conn = conn

	if err := db.runMigrations(); err != nil {
		return err
	}

	return nil
}

func (db *DB) runMigrations() error {
	log.Printf("running migrations")
	url := db.URL()

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "could not create migration source: %s", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "could not create migration from source: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return talenthub.Errorf(talenthub.EINTERNAL, "failed to run migrations: %s", err)
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *pgx.TxOptions) (*Tx, error) {
	if opts == nil {
		opts = &pgx.TxOptions{}
	}
	tx, err := db.conn.BeginTx(ctx, *opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		now: time.Now(),
	}, nil
}

type Tx struct {
	pgx.Tx
	now time.Time
}
