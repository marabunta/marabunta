package marabunta

// Flags available command flags
type Flags struct {
	CA         string
	Configfile string
	Crt        string
	GRPC       uint
	HTTP       uint
	Key        string
	Mysql      string
	Redis      string
	Version    bool
}
