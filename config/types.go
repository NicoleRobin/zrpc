package config

type HostConfig struct{}

type LogConfig struct {
	Directory string `yaml:"directory"`
}

type AppConfig struct {
	AppName    string    `yaml:"app_name"`
	EnableHTTP bool      `yaml:"enable_http"`
	Address    string    `yaml:"address"`
	Log        LogConfig `yaml:"log"`
}

type ClientConfig struct {
	AppName     string `yaml:"app_name"`
	ServiceName string `yaml:"service_name"`
	DialTimeout int    `yaml:"dial_timeout"`
	Provider    string `yaml:"provider"`
}
