package yaml

import "fmt"

// Postgres holds the configuration for the PostgreSQL database.
type Postgres struct {
	host     string `validate:"required" yaml:"host"`
	port     string `validate:"required" yaml:"port"`
	user     string `validate:"required" yaml:"user"`
	password string `validate:"required" yaml:"password"`
	dbName   string `validate:"required" yaml:"dbname"`
	sslMode  string `validate:"required" yaml:"sslmode"`
}

// DSN returns the Data Source Name for connecting to PostgreSQL.
func (p Postgres) DSN() string {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host,
		p.port,
		p.user,
		p.password,
		p.dbName,
		p.sslMode,
	)

	return connStr
}
