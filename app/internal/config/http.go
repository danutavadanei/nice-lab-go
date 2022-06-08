package config

import (
	"github.com/spf13/viper"
	"time"
)

// HTTPServerConfig stores the configuration server server
type HTTPServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewHTTPServerConfig returns a new HTTPServerConfig
func NewHTTPServerConfig(v *viper.Viper) HTTPServerConfig {
	v.SetDefault("HTTP_ADDR", ":8080")

	return HTTPServerConfig{
		Addr:         v.GetString("HTTP_ADDR"),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 120,
	}
}
