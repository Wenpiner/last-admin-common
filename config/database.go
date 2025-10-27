package config

import "fmt"

type DatabaseConfig struct {
	Host         string `json:",env=DATABASE_HOST"`                                                   // 数据库地址
	Port         int    `json:",env=DATABASE_PORT"`                                                   // 数据库端口
	Username     string `json:",env=DATABASE_USERNAME"`                                               // 数据库用户名
	Password     string `json:",env=DATABASE_PASSWORD"`                                               // 数据库密码
	DatabaseName string `json:",env=DATABASE_NAME"`                                                   // 数据库名称
	SSLMode      string `json:",env=DATABASE_SSL_MODE"`                                               // SSL 模式
	DBType       string `json:",env=DATABASE_TYPE,default=postgres,options=[mysql,postgres,sqlite3]"` // 数据库类型
	MaxOpenConns int    `json:",optional,env=DATABASE_MAX_OPEN_CONNS,default=100"`
	MaxIdleConns int    `json:",optional,env=DATABASE_MAX_IDLE_CONNS,default=10"`
	CacheTime    int    `json:",optional,env=DATABASE_CACHE_TIME,default=300"`
	DBPath       string `json:",optional,env=DATABASE_PATH"`
}

func (c *DatabaseConfig) GetDSN() string {
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

func (c *DatabaseConfig) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
}

func (c *DatabaseConfig) PostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.DatabaseName, c.SSLMode)
}

func (c *DatabaseConfig) SqliteDSN() string {
	return fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1", c.DBPath)
}
