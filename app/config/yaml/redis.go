package yaml

import (
	"net"
	"time"
)

// Redis represents the configuration for Redis.
type Redis struct {
	Host              string        `validate:"required" yaml:"host"`
	Port              string        `validate:"required" yaml:"port"`
	ConnectionTimeout time.Duration `validate:"required" yaml:"connection_timeout"`
	MaxIdle           int           `validate:"required" yaml:"max_idle"`
	IdleTimeout       time.Duration `validate:"required" yaml:"idle_timeout"`
}

// Address returns the Redis server address.
func (r Redis) Address() string {
	return net.JoinHostPort(r.Host, r.Port)
}
