package config

import (
	"time"
)

type Config struct {
	URL      string        `json:"url,omitempty" mapstructure:"url,omitempty"`
	APIURL   string        `json:"apiurl,omitempty" mapstructure:"apiurl,omitempty"`
	Username string        `json:"username,omitempty" mapstructure:"username,omitempty"`
	Password string        `json:"password,omitempty" mapstructure:"password,omitempty"`
	Token    string        `json:"token,omitempty" mapstructure:"token,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty" mapstructure:"timeout,omitempty"`
	Debug    bool          `json:"debug,omitempty" mapstructure:"debug,omitempty"`
	Me       string        `json:"me,omitempty" mapstructure:"me,omitempty"`
}
