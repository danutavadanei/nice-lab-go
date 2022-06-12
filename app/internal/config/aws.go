package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
)

func NewAWSConfig(v *viper.Viper) (cfg aws.Config, err error) {
	v.SetDefault("AWS_REGION", "us-east-1")
	v.SetDefault("AWS_ACCESS_KEY_ID", "key")
	v.SetDefault("AWS_SECRET_ACCESS_KEY", "secret")
	v.SetDefault("AWS_S3_DEFAULT_BUCKET", "default-bucket")

	cfg, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				v.GetString("AWS_ACCESS_KEY_ID"),
				v.GetString("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		),
		config.WithRegion(v.GetString("AWS_REGION")),
	)

	return
}
