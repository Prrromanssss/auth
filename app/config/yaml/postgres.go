package yaml

import "fmt"

// Postgres holds the configuration for the PostgreSQL database.
type Postgres struct {
	Host     string `validate:"required" yaml:"host"`
	Port     string `validate:"required" yaml:"port"`
	User     string `validate:"required" yaml:"user"`
	Password string `validate:"required" yaml:"password"`
	DBName   string `validate:"required" yaml:"dbname"`
	SSLMode  string `validate:"required" yaml:"sslmode"`
}

// DSN returns the Data Source Name for connecting to PostgreSQL.
func (p Postgres) DSN() string {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DBName,
		p.SSLMode,
	)

	return connStr
}
