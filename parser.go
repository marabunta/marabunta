package marabunta

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// Parser command line options and configuration file
type Parse struct {
	Flags
}

// Parse parse the command line flags
func (p *Parse) Parse(fs *flag.FlagSet) (*Flags, error) {
	fs.BoolVar(&p.Flags.Version, "v", false, "Print version")
	fs.StringVar(&p.Flags.Configfile, "c", "", "`marabunta.yml` configuration file")
	fs.StringVar(&p.Flags.Mysql, "mysql", "", "MySQL `DSN` username:password@address:port/dbname")
	fs.StringVar(&p.Flags.Redis, "redis", "", "Redis `host:port`")
	fs.UintVar(&p.Flags.GRPC, "grpc", 1415, "Listen on gRPC `port` default 1415")
	fs.UintVar(&p.Flags.HTTP, "http", 8000, "Listen on HTTP `port` default 8000")
	fs.StringVar(&p.Flags.CA, "tls.ca", "", "Path to TLS Certificate Authority (`CA`)")
	fs.StringVar(&p.Flags.Crt, "tls.crt", "", "Path to TLS `certificate`")
	fs.StringVar(&p.Flags.Key, "tls.key", "", "Path to TLS `private key`")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &p.Flags, nil
}

func (p *Parse) parseYml(file string) (*Config, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// set defaults
	var cfg = Config{
		HTTPPort: 8000,
		GRPCPort: 1415,
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse YAML file %q %s", file, err)
	}
	return &cfg, nil
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
func (p *Parse) ParseArgs(fs *flag.FlagSet) (cfg *Config, err error) {
	flags, err := p.Parse(fs)
	if err != nil {
		return
	}

	// if -v
	if flags.Version {
		return
	}

	// if -c
	if flags.Configfile != "" {
		if !isFile(flags.Configfile) {
			err = fmt.Errorf("cannot read file: %q, use (\"%s -h\") for help", flags.Configfile, os.Args[0])
			return
		}

		// parse the `run.yml` file
		cfg, err = p.parseYml(flags.Configfile)
		if err != nil {
			return
		}

		return
	}

	// if no args
	if len(fs.Args()) < 1 {
		err = fmt.Errorf("missing options, use (\"%s -h\") for help", os.Args[0])
		return
	}

	// create new cfg if not using -c
	return
}
