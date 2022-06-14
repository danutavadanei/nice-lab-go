package config

import (
	"github.com/spf13/viper"
)

type GatewayConfig struct {
	AuthUrl     string
	PipelineUrl string
}

func NewGatewayConfig(v *viper.Viper) GatewayConfig {
	v.SetDefault("AUTH_SERVICE_URL", "http://auth:8080")
	v.SetDefault("PIPELINE_SERVICE_URL", "http://pipeline:8080")

	return GatewayConfig{
		AuthUrl:     v.GetString("AUTH_SERVICE_URL"),
		PipelineUrl: v.GetString("PIPELINE_SERVICE_URL"),
	}
}
