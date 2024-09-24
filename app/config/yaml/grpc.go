package yaml

import "net"

// GRPCServer holds the configuration for the gRPC server.
type GRPCServer struct {
	Host string `validate:"required" yaml:"host"`
	Port string `validate:"required" yaml:"port"`
}

// Address returns the server address.
func (s GRPCServer) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}
