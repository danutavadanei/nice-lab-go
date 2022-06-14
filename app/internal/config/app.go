package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AWSConfig        *aws.Config
	HTTPServerConfig HTTPServerConfig
	MySQLConfig      mysql.Config
	GatewayConfig    GatewayConfig
}

func NewAppConfig(v *viper.Viper) AppConfig {
	return AppConfig{
		AWSConfig:        NewAWSConfig(v),
		HTTPServerConfig: NewHTTPServerConfig(v),
		MySQLConfig:      NewMySQLConfig(v),
		GatewayConfig:    NewGatewayConfig(v),
	}
}
