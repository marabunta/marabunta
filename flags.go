package marabunta

// Flags available command flags
type Flags struct {
	Configfile string
	GRPC       int
	HTTP       int
	Mysql      string
	Redis      string
	TLSCACrt   string
	TLSCAKey   string
	TLSCrt     string
	TLSKey     string
	Version    bool
}
