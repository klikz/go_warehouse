package apiserver

type Config struct {
	s_address   string `toml:"saddress"`
	logLevel    string `toml:"logLevel"`
	DatabaseURL string `toml:"database_url"`
	SecretKey   string `toml:"secretkey"`
}

func NewConfig() *Config {
	return &Config{
		s_address: ":7777",
		logLevel:  "debug",
	}
}
