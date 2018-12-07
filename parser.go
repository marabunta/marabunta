package marabunta

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// Parse command line options and configuration file
type Parse struct {
	Flags
}

// Parse parse the command line flags
func (p *Parse) Parse(fs *flag.FlagSet) (*Flags, error) {
	fs.BoolVar(&p.Flags.Version, "v", false, "Print version")
	fs.StringVar(&p.Flags.Configfile, "c", "", "`marabunta.yml` configuration file")
	fs.StringVar(&p.Flags.Mysql, "mysql", "", "MySQL `DSN` username:password@address:port/dbname")
	fs.StringVar(&p.Flags.Redis, "redis", "", "Redis `host:port` (default 127.0.0.1:6379)")
	fs.IntVar(&p.Flags.GRPC, "grpc", 1415, "Listen on gRPC `port` (default 1415)")
	fs.IntVar(&p.Flags.HTTP, "http", 8000, "Listen on HTTP `port` (default 8000)")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &p.Flags, nil
}

func (p *Parse) parseYml(file string, cfg *Config) (*Config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse YAML file %q %s", file, err)
	}
	return cfg, nil
}

// Usage prints to standard error a usage message
func (p *Parse) Usage(fs *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...]\n\n", os.Args[0])
		var flags []string
		fs.VisitAll(func(f *flag.Flag) {
			flags = append(flags, f.Name)
		})
		for _, v := range flags {
			f := fs.Lookup(v)
			s := fmt.Sprintf("  -%s", f.Name)
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}
			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}
			s += usage
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
	}
}

// ParseArgs parse command arguments
func (p *Parse) ParseArgs(fs *flag.FlagSet) (*Config, error) {
	flags, err := p.Parse(fs)
	if err != nil {
		return nil, err
	}

	// if -v
	if flags.Version {
		return nil, nil
	}

	var (
		needCertificate bool
		cfg             = &Config{
			HTTPPort: 8000,
			GRPCPort: 1415,
		}
	)

	// if -c
	if flags.Configfile != "" {
		if !isFile(flags.Configfile) {
			return nil, fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.Configfile, os.Args[0])
		}

		// parse the `run.yml` file
		cfg, err := p.parseYml(flags.Configfile, cfg)
		if err != nil {
			return nil, err
		}

		// Home
		if cfg.Home == "" {
			home, err := GetHome()
			if err != nil {
				return nil, err
			}
			cfg.Home = home
		}

		// TLS CA crt
		if cfg.TLS.CACrt != "" {
			if !isFile(cfg.TLS.CACrt) {
				return nil, fmt.Errorf("cannot read TLS CA crt file: %q, use (\"%s -h\") for help", cfg.TLS.CACrt, os.Args[0])
			}
		} else {
			cfg.TLS.CACrt = filepath.Join(cfg.Home, "CA.crt")
			needCertificate = true
		}

		// TLS CA key
		if cfg.TLS.CAKey != "" {
			if !isFile(cfg.TLS.CAKey) {
				return nil, fmt.Errorf("cannot read TLS CA key file: %q, use (\"%s -h\") for help", cfg.TLS.CAKey, os.Args[0])
			}
		} else {
			cfg.TLS.CAKey = filepath.Join(cfg.Home, "CA.key")
			needCertificate = true
		}

		// TLS certificate
		if cfg.TLS.Crt != "" {
			if !isFile(cfg.TLS.Crt) {
				return nil, fmt.Errorf("cannot read TLS crt file: %q, use (\"%s -h\") for help", cfg.TLS.Crt, os.Args[0])
			}
		} else {
			cfg.TLS.Crt = filepath.Join(cfg.Home, "http.crt")
			needCertificate = true
		}

		// TLS KEY
		if cfg.TLS.Key != "" {
			if !isFile(cfg.TLS.Key) {
				return nil, fmt.Errorf("cannot read TLS Key file: %q, use (\"%s -h\") for help", cfg.TLS.Key, os.Args[0])
			}
		} else {
			cfg.TLS.Key = filepath.Join(cfg.Home, "http.key")
			needCertificate = true
		}

		if needCertificate {
			if err := createCertificates(cfg); err != nil {
				return nil, err
			}
		}

		return cfg, nil
	}

	if fs.NFlag() < 1 {
		return nil, fmt.Errorf("missing options, use (\"%s -h\") for help", os.Args[0])
	}

	home, err := GetHome()
	if err != nil {
		return nil, err
	}

	cfg.Home = home

	if flags.GRPC != 0 {
		cfg.GRPCPort = flags.GRPC
	}

	if flags.HTTP != 0 {
		cfg.HTTPPort = flags.HTTP
	}

	if flags.Mysql != "" {
		// TODO parse DSN
	} else {
		return nil, fmt.Errorf("missing MySQL DSN, use (\"%s -h\") for help", os.Args[0])
	}

	if flags.Redis != "" {
		// TODO parse redis
	} else {
		cfg.Redis = Redis{"127.0.0.1", 6379}
	}

	// check if certs already exists
	cfg.TLS.CACrt = filepath.Join(cfg.Home, "CA.crt")
	if !isFile(cfg.TLS.CACrt) {
		needCertificate = true
	}
	cfg.TLS.CAKey = filepath.Join(cfg.Home, "CA.key")
	if !isFile(cfg.TLS.CAKey) {
		needCertificate = true
	}
	cfg.TLS.Key = filepath.Join(cfg.Home, "http.key")
	if !isFile(cfg.TLS.Key) {
		needCertificate = true
	}
	cfg.TLS.Crt = filepath.Join(cfg.Home, "http.crt")
	if !isFile(cfg.TLS.Crt) {
		needCertificate = true
	}

	if needCertificate {
		if err := createCertificates(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
