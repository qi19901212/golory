// golory is ALL IN ONE package for go software
// development with best practice usages support
package golory

import (
	"encoding/json"
	"fmt"
	"github.com/1pb-club/golory/log"
	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

var (
	gly *golory
)

// golory struct is used to hold all data
type golory struct {
	cfg        *goloryConfig
	components *handler
	booted     bool
}

// goloryConfig is used to store golory configurations
type goloryConfig struct {
	Debug   bool
	Loggers map[string]log.CommonCfg
}

func init() {
	gly = &golory{
		booted:     false,
		cfg:        &goloryConfig{},
		components: newHandler(),
	}
}

// Initiate components from configuration file or binary content.
// Toml, Json, Yaml supported.
func Boot(cfg interface{}) error {
	if gly.booted {
		// do clear stuff
		gly.booted = false
	}
	switch cfg.(type) {
	case string:
		if err := parseFile(cfg.(string)); err != nil {
			return err
		}
	case []byte:
		if err := parseBytes(cfg.([]byte)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("only string or []byte supported, %s", cfg)
	}
	// do initiation
	gly.init()
	gly.booted = true
	return nil
}

// Initate golory components from file
func parseFile(path string) error {
	// read file to []byte
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return parseBytes(b)
}

// Initiate golory components from binary content
func parseBytes(b []byte) error {
	if err := parseCfg(b); err != nil {
		return err
	}
	return nil
}

// Do parse config.
// It will try serveral formats one by one
func parseCfg(b []byte) error {
	// try file formats
	var err error
	if err = toml.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	e := wrap(GLY_PARSE_CFG_ERROR, err)
	if err = yaml.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	e = wrap(e, err)
	if err = json.Unmarshal(b, &gly.cfg); err == nil {
		return nil
	}
	return wrap(e, err)
}

// Init all components
func (g *golory) init() {
	g.initLog()
}

// Init log component
func (g *golory) initLog() {
	if g.cfg.Loggers == nil {
		return
	}
	for key, cfg := range g.cfg.Loggers {
		// TODO log.Boot should return error when Boot failed
		obj := log.Boot(cfg)
		g.components.setLogger(key, obj)
	}
}
