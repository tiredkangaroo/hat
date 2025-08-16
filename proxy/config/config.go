package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Addr string `toml:"addr"`

	MITM struct {
		CertificateFile     string `toml:"certificate_file"`
		KeyFile             string `toml:"key_file"`
		CertificateLifetime int64  `toml:"certificate_lifetime"`
	} `toml:"mitm"`
}

var DefaultConfig = &Configuration{}

func (c *Configuration) Init() error {
	// get the configuration file
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get user config dir: %w", err)
	}
	configFilename := filepath.Join(configDir, "hat", "config.toml")
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", configFilename)
	}
	// read the configuration file
	configFile, err := os.OpenFile(configFilename, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	defer configFile.Close()

	md, err := toml.NewDecoder(configFile).Decode(c)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}
	if len(md.Undecoded()) > 0 {
		return fmt.Errorf("config file contains undecoded fields: %v", md.Undecoded())
	}

	if c.Addr == "" || c.MITM.CertificateFile == "" || c.MITM.KeyFile == "" {
		return fmt.Errorf("config file is missing some required fields")
	}

	return nil
}

func Init() error {
	return DefaultConfig.Init()
}
