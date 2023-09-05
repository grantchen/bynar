package config

type sendgridConfig struct {
	ApiKey       string
	FromName     string
	FromAddress  string
	ToName       string
	MagicLinkUrl string
}

const sendgridApiKey = "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"
const fromName = "Bynar"
const fromAddress = "wilson_wz@163.com"
const toName = ""
const magicLinkUrl = "http://xxx.com/"

// Get sendgrid config
func GetSendgridConfig() sendgridConfig {
	config := sendgridConfig{
		ApiKey:       sendgridApiKey,
		FromName:     fromName,
		FromAddress:  fromAddress,
		ToName:       toName,
		MagicLinkUrl: magicLinkUrl,
	}
	return config
}
