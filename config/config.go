package config

import (
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"os"
)

const Location = "~/.config/jump"

type Config struct {
	Tunnels map[string]Tunnel `toml:"tunnels"`
	Remotes map[string]Remote `toml:"remotes"`
}

type Tunnel struct {
	LocalPort  int16  `toml:"local"`
	RemotePort int16  `toml:"remote"`
	Addr       string `toml:"address"`
}

type Remote struct {
	Addr      string  `toml:"address"`
	Port      int16   `toml:"port"`
	User      *string `toml:"user"`
	IdentFile *string `toml:"identity"`
}

func (c Config) GetRemote(name string) *Remote {
	if c.Remotes == nil {
		return nil
	}

	if r, exists := c.Remotes[name]; exists {
		return &r
	} else {
		return nil
	}
}

func (c *Config) AddRemote(name string, remote Remote) {
	if c.Remotes == nil {
		c.Remotes = map[string]Remote{}
	}

	c.Remotes[name] = remote
}

func (c Config) GetTunnel(name string) *Tunnel {
	if c.Tunnels == nil {
		return nil
	}

	if t, exists := c.Tunnels[name]; exists {
		return &t
	} else {
		return nil
	}

}

func (c *Config) AddTunnel(name string, tunnel Tunnel) {
	if c.Tunnels == nil {
		c.Tunnels = map[string]Tunnel{}
	}

	c.Tunnels[name] = tunnel
}

func Read(cfgFile string) (Config, error) {
	loc, _ := homedir.Expand(cfgFile)
	c := Config{}
	file, err := os.Open(loc)
	if err != nil {
		return c, nil
	}

	if _, err := toml.DecodeReader(file, &c); err != nil {
		return c, errors.Wrap(err, "failed to read config file")
	}
	return c, nil
}

func Write(cfgFile string, cfg Config) (Config, error) {
	loc, _ := homedir.Expand(cfgFile)
	file, err := os.OpenFile(loc, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to open config file for writing")
	}

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		return cfg, errors.Wrapf(err, "failed to write config to file %v", cfgFile)
	}
	return Read(cfgFile)
}
