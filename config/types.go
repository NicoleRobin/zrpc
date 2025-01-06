package config

type HostConfig struct{}

type LogConfig struct {
}

type AppConfig struct {
	AppName    string    `yaml:"app_name"`
	EnableHTTP bool      `yaml:"enable_http"`
	Addr       string    `yaml:"addr"`
	Log        LogConfig `yaml:"log"`
}
