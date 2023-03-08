package config

// read from toml format config file
type ServerConfig struct {
	ServerPort     string
	ServerPath     string
	ServerMode     string
	MaxMessageSize int64

	PGSQL struct {
		DBName   string
		Username string
		Password string
		Host     string
		Port     string
	}
}
