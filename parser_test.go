package marabunta

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParseHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-h"}
	p := &Parse{}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Usage = p.Usage(fs)
	// Error output buffer
	buf := bytes.NewBuffer([]byte{})
	fs.SetOutput(buf)
	_, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	os.Stderr = w
	_, err = p.Parse(fs)
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestParseDefault(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", ""}
	p := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	flags, err := p.Parse(fs)
	if err != nil {
		t.Error(err)
	}
	if helpCalled {
		t.Error("help called for regular flag")
	}
	expect(t, "", flags.Configfile)
	expect(t, "", flags.Mysql)
	expect(t, "", flags.Redis)
	expect(t, "", flags.TLSCACrt)
	expect(t, "", flags.TLSCAKey)
	expect(t, "", flags.TLSCrt)
	expect(t, "", flags.TLSKey)
	expect(t, int(1415), flags.GRPC)
	expect(t, int(8000), flags.HTTP)
	expect(t, false, flags.Version)
}

func TestParseFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	var flagTest = []struct {
		flag     []string
		name     string
		expected interface{}
	}{
		{[]string{"cmd", "-v"}, "Version", true},
		{[]string{"cmd", "-c", "marabunta.yml"}, "Configfile", "marabunta.yml"},
		{[]string{"cmd", "-mysql", "username:password@host:port/database"}, "mysql", "username:password@host:port/database"},
		{[]string{"cmd", "-redis", "host:port"}, "redis", "host:port"},
		{[]string{"cmd", "-grpc", "1415"}, "gRPC", "1415"},
		{[]string{"cmd", "-http", "8000"}, "http", "8000"},
	}

	var helpCalled = false
	for _, f := range flagTest {
		os.Args = f.flag
		p := &Parse{}
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.Usage = func() { helpCalled = true }
		flags, err := p.Parse(fs)
		if err != nil {
			t.Error(err)
		}
		if helpCalled {
			t.Error("help called for regular flag")
			helpCalled = false // reset for next test
		}
		refValue := reflect.ValueOf(flags).Elem().FieldByName(f.name)
		switch refValue.Kind() {
		case reflect.Bool:
			expect(t, f.expected, refValue.Bool())
		case reflect.String:
			expect(t, f.expected, refValue.String())
		case reflect.Int:
			expect(t, f.expected, int(refValue.Int()))
		case reflect.Uint:
			expect(t, uint(f.expected.(int)), uint(refValue.Uint()))
		}
	}
}

func TestParseArgsHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-h"}
	parser := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("TestParseArgsHelp", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	parser.ParseArgs(fs)
	if !helpCalled {
		t.Fatal("help was not called")
	}
}

func TestParseArgsVersion(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-v"}
	parser := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("TestParseArgsVersion", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	parser.ParseArgs(fs)
	if helpCalled {
		t.Error("help called for regular flag")
	}
}

func TestParseArgsVersion2(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-v", "-c", "xyz"}
	parser := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("TestParseArgsVersion2", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	parser.ParseArgs(fs)
	if helpCalled {
		t.Error("help called for regular flag")
	}
}

func TestParseArgsNoargs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}
	parser := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("TestParseArgsNoargs", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	_, err := parser.ParseArgs(fs)
	if helpCalled {
		t.Error("help called for regular flag")
	}
	if err == nil {
		t.Error("Expecting error")
	}
}

func TestParseArgsTable(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	var flagTest = []struct {
		name        string
		flag        []string
		expectError bool
	}{
		{"version", []string{"cmd", "-v"}, false},
		{"1.yml", []string{"cmd", "-c", "marabunta.yml"}, true},
		{"2.yml", []string{"cmd", "-c", "example/marabunta.yml", "cmd"}, false},
		// todo fix parse dns
		{"dsn", []string{"cmd", "-mysql", "dsn"}, false},
		{"redis", []string{"cmd", "-redis", "host:port"}, true},
	}
	var helpCalled = false
	for _, tc := range flagTest {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = tc.flag
			parser := &Parse{}
			fs := flag.NewFlagSet("TestParseArgsTable", flag.ContinueOnError)
			fs.Usage = func() { helpCalled = true }
			_, err := parser.ParseArgs(fs)
			if tc.expectError {
				if err == nil {
					t.Error("Expecting error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
			if helpCalled {
				t.Error("help called for regular flag")
				helpCalled = false // reset for next test
			}
		})
	}
}

func TestParseYamlCmd(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestParseYamlCmd")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())
	yaml := []byte(`
http_port: 8000
grpc_port: 1415
mysql:
  host: localhost
  port: 3306
  database: marabunta
  username: root
  password: example
redis:
  host: localhost
  port: 6379
tls:
  crt: certs/server.crt
  key: certs/server.key
  ca: certs/CA.crt`)
	err = ioutil.WriteFile(tmpfile.Name(), yaml, 0644)
	if err != nil {
		t.Error(err)
	}
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-c", tmpfile.Name()}
	parser := &Parse{}
	var helpCalled = false
	fs := flag.NewFlagSet("TestParseArgsYaml", flag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }
	_, err = parser.ParseArgs(fs)
	if helpCalled {
		t.Error("help called for regular flag")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestParseParseYmlioutil(t *testing.T) {
	p := &Parse{}
	c := &Config{}
	if _, err := p.parseYml("/dev/null/non-existent", c); err == nil {
		t.Error("Expecting error")
	}
}

func TestParseBadYaml(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestParseBadYaml")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile.Name())
	yaml := []byte(`
grpc - command
http: 10`)
	err = ioutil.WriteFile(tmpfile.Name(), yaml, 0644)
	if err != nil {
		t.Error(err)
	}
	p := &Parse{}
	c := &Config{}
	if _, err := p.parseYml(tmpfile.Name(), c); err == nil {
		t.Error("Expecting error")
	}
}
