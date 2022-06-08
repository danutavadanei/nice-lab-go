package aws_test

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/danutavadanei/localstack-go-playground/internal/adapters/aws"
	"github.com/danutavadanei/localstack-go-playground/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var awsClient aws.Client

func TestMain(m *testing.M) {
	v := viper.New()
	v.AutomaticEnv()
	cfg := config.NewAppConfig(v)
	awsClient = aws.NewClient(cfg.AWSConfig)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestClient_ListBuckets(t *testing.T) {
	t.Parallel()

	bytes, err := awsClient.ListBuckets(context.Background())

	assert.NoError(t, err)

	var buckets *s3.ListBucketsOutput

	err = json.Unmarshal(bytes, &buckets)

	assert.NoError(t, err)

	assert.Equal(t, 3, len(buckets.Buckets))
}
