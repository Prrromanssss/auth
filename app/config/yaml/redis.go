package yaml

import (
	"net"
	"time"
)

// Redis represents the configuration for Redis.
type Redis struct {
	host              string        `validate:"required" yaml:"host"`
	port              string        `validate:"required" yaml:"port"`
	connectionTimeout time.Duration `validate:"required" yaml:"connection_timeout"`
	maxIdle           int           `validate:"required" yaml:"max_idle"`
	idleTimeout       time.Duration `validate:"required" yaml:"idle_timeout"`
}

// Address returns the Redis server address.
func (r Redis) Address() string {
	return net.JoinHostPort(r.host, r.port)
}

// ConnectionTimeout returns the connection timeout duration.
func (r Redis) ConnectionTimeout() time.Duration {
	return r.connectionTimeout
}

// MaxIdle returns the maximum number of idle connections.
func (r Redis) MaxIdle() int {
	return r.maxIdle
}

// IdleTimeout returns the idle timeout duration.
func (r Redis) IdleTimeout() time.Duration {
	return r.idleTimeout
}
