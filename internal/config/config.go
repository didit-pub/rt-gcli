package config

import (
	"time"
)

type Config struct {
	URL      string        `mapstructure:"url"`
	APIURL   string        `mapstructure:"apiurl"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Token    string        `mapstructure:"token"`
	Timeout  time.Duration `mapstructure:"timeout"`
	Debug    bool          `mapstructure:"debug"`
	Me       string        `mapstructure:"me"`
}
