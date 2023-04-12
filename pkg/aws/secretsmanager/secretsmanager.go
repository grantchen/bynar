package secretsmanager

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	VersionStage    string
	MaxCacheSize    int
	CacheItemTTL    int64
}

type awsSecretsManager struct {
	cache *secretcache.Cache
}

func NewAWSSecretsManager(cfg AWSConfig) (SecretsManager, error) {
	client, err := createAWSSecretsManagerClient(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.Region)
	if err != nil {
		return nil, err
	}

	config := secretcache.CacheConfig{
		VersionStage: cfg.VersionStage,
		MaxCacheSize: cfg.MaxCacheSize,
		CacheItemTTL: cfg.CacheItemTTL,
	}

	cache, err := secretcache.New(
		func(c *secretcache.Cache) { c.Client = client },
		func(c *secretcache.Cache) { c.CacheConfig = config },
	)

	if err != nil {
		return nil, err
	}

	return &awsSecretsManager{
		cache: cache,
	}, nil
}

func (rx *awsSecretsManager) GetString(key string) (string, error) {
	return rx.cache.GetSecretString(key)
}

func (rx *awsSecretsManager) GetJSONData(key string, out interface{}) error {
	secretString, err := rx.cache.GetSecretString(key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(secretString), &out)
}

func createAWSSecretsManagerClient(accessKeyID, secretAccessKey, region string) (secretsmanageriface.SecretsManagerAPI, error) {
	session, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion(region)

	if accessKeyID != "" && secretAccessKey != "" {
		cfg.WithCredentials(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""))
	}

	return secretsmanager.New(session, cfg), nil
}
