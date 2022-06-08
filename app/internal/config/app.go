package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AWSConfig        *aws.Config
	HTTPServerConfig HTTPServerConfig
}

// NewAppConfig creates a new AppConfig
func NewAppConfig(v *viper.Viper) AppConfig {
	return AppConfig{
		AWSConfig:        NewAWSConfig(v),
		HTTPServerConfig: NewHTTPServerConfig(v),
	}
}
