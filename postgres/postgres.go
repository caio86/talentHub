package postgres

import (
	"context"
	"fmt"

	talenthub "github.com/caio86/talentHub"
	"github.com/jackc/pgx/v5"
)

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
	conn *pgx.Conn
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

	conn, err := pgx.Connect(context.Background(), db.URL())
	if err != nil {
		return talenthub.Errorf(talenthub.EINTERNAL, "%s", err)
	}

	db.conn = conn

	return nil
}
