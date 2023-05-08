package secretsmanager

type SecretsManager interface {
	GetString(key string) (string, error)
	GetJSONData(key string, out interface{}) error
}
