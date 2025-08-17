package baseserver

type Configuration struct {
	Services ServicesConfiguration `json:"services" yaml:"services"`
}

type ServicesConfiguration struct {
	GRPC *ServerConfiguration `json:"grpc,omitempty" yaml:"grpc,omitempty"`
	HTTP *ServerConfiguration `json:"http,omitempty" yaml:"http,omitempty"`
}

type ServerConfiguration struct {
	Address string            `json:"address" yaml:"address"`
	TLS     *TLSConfiguration `json:"tls,omitempty" yaml:"tls,omitempty"`
}

// GetAddress returns the configured address or an empty string of s is nil
func (s *ServerConfiguration) GetAddress() string {
	if s == nil {
		return ""
	}
	return s.Address
}

type TLSConfiguration struct {
	CAPath   string `json:"caPath" yaml:"caPath"`
	CertPath string `json:"certPath" yaml:"certPath"`
	KeyPath  string `json:"keyPath" yaml:"keyPath"`
}
