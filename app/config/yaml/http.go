package yaml

import "net"

// HTTPServer holds the configuration for the HTTP server.
type HTTPServer struct {
	Host              string `validate:"required" yaml:"host"`
	Port              string `validate:"required" yaml:"port"`
	ReadHeaderTimeout int64  `validate:"required" yaml:"read_header_timeout"`
}

// Address returns the server address.
func (s HTTPServer) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}
