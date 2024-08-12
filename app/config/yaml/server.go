package yaml

import "net"

// Server holds the configuration for the gRPC/HTTP server.
type Server struct {
	Host string `validate:"required" yaml:"host"`
	Port string `validate:"required" yaml:"port"`
}

// Address returns the server address.
func (s Server) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}
