package marabunta

// Config yaml/command line configuration
type Config struct {
	HTTPPort int `yaml:"http_port"`
	GRPCPort int `yaml:"grpc_port"`
	MySQL    `yaml:"mysql"`
	Redis    `yaml:"redis"`
	TLS      `yaml:"tls"`
}

// MySQL configuration options
type MySQL struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Redis configuration options
type Redis struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

// TLS configuration options
type TLS struct {
	Crt   string `yaml:"crt"`
	Key   string `yaml:"key"`
	CACrt string `yaml:"ca_crt"`
	CAKey string `yaml:"ca_key"`
}
