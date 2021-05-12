package connect

type Config struct {
	AppName       string `toml:"app_name"`
	AppId         string `toml:"app_id"`
	KeyId         string `toml:"key_id,omitempty"`
	AppSecret     string `toml:"app_secret"`
	AccessToken   string `toml:"access_token"`
	OAuthCallback string `toml:"oauth_callback"`
	WebhookURL    string `toml:"webhook_url,omitempty"`
	Scope         string `toml:"scope,omitempty"`
}

type Configs struct{}

// Configure loads the connect configs from config file
func Configure(confs Configs) {
}
