package marabunta

// Flags available command flags
type Flags struct {
	Configfile string
	GRPC       uint
	HTTP       uint
	Mysql      string
	Redis      string
	TLSCA      string
	TLSCrt     string
	TLSKey     string
	Version    bool
}
