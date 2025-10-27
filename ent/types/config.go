package types

import (
	"context"
	"database/sql"
	esql "entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"time"
)

type DatabaseConfig struct {
	Host         string `json:",env=DATABASE_HOST"`
	Port         int    `json:",env=DATABASE_PORT"`
	Username     string `json:",env=DATABASE_USERNAME"`
	Password     string `json:",env=DATABASE_PASSWORD"`
	DatabaseName string `json:",env=DATABASE_NAME"`
	SSLMode      string `json:",env=DATABASE_SSL_MODE"`
	DBType       string `json:",env=DATABASE_TYPE,default=postgres,options=[mysql,postgres,sqlite3]"`
	MaxOpenConns int    `json:",optional,env=DATABASE_MAX_OPEN_CONNS,default=100"`
	MaxIdleConns int    `json:",optional,env=DATABASE_MAX_IDLE_CONNS,default=10"`
	CacheTime    int    `json:",optional,env=DATABASE_CACHE_TIME,default=300"`
	DBPath       string `json:",optional,env=DATABASE_PATH"`
	Config       string `json:",optional,env=DATABASE_CONFIG"`
}

// GetDSN returns DSN according to the database type.
func (c DatabaseConfig) GetDSN() string {
	switch c.DBType {
	case "mysql":
		return c.MysqlDSN()
	case "postgres":
		return c.PostgresDSN()
	case "sqlite3":
		return c.SqliteDSN()
	default:
		return "mysql"
	}
}

// MysqlDSN returns mysql DSN.
func (c DatabaseConfig) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True%s", c.Username, c.Password, c.Host, c.Port, c.DatabaseName, c.Config)
}

// PostgresDSN returns Postgres DSN.
func (c DatabaseConfig) PostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s%s", c.Username, c.Password, c.Host, c.Port, c.DatabaseName,
		c.SSLMode, c.Config)
}

// SqliteDSN returns Sqlite DSN.
func (c DatabaseConfig) SqliteDSN() string {
	if c.DBPath == "" {
		logx.Must(errors.New("the database file path cannot be empty"))
	}

	if _, err := os.Stat(c.DBPath); os.IsNotExist(err) {
		f, err := os.OpenFile(c.DBPath, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			logx.Must(fmt.Errorf("failed to create SQLite database file %q", c.DBPath))
		}
		if err := f.Close(); err != nil {
			logx.Must(fmt.Errorf("failed to create SQLite database file %q", c.DBPath))
		}
	} else {
		if err := os.Chmod(c.DBPath, 0660); err != nil {
			logx.Must(fmt.Errorf("unable to set permission code on %s: %v", c.DBPath, err))
		}
	}

	return fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1%s", c.DBPath, c.Config)
}

// NewNoCacheDriver 新建一个不带缓存的驱动
func (c DatabaseConfig) NewNoCacheDriver() *esql.Driver {
	db, err := sql.Open(c.DBType, c.GetDSN())
	logx.Must(err)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = db.PingContext(ctx)
	logx.Must(err)

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	return esql.OpenDB(c.DBType, db)
}
