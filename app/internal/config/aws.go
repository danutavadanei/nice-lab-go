package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/viper"
)

func NewAWSConfig(v *viper.Viper) *aws.Config {
	v.SetDefault("AWS_ENDPOINT", "http://localhost:4566")
	v.SetDefault("AWS_REGION", "us-west-2")
	v.SetDefault("AWS_S3_DEFAULT_BUCKET", "default-bucket")

	awsCfg := aws.NewConfig()
	awsCfg.WithEndpoint(v.GetString("AWS_ENDPOINT"))
	awsCfg.WithS3ForcePathStyle(true)
	awsCfg.WithRegion(v.GetString("AWS_REGION"))
	awsCfg.WithCredentialsChainVerboseErrors(true)
	awsCfg.WithCredentials(credentials.NewStaticCredentials("test", "test", "test"))

	return awsCfg
}
