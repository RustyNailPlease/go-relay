package config

// read from toml format config file
type ServerConfig struct {
	ServerPort     string
	ServerPath     string
	ServerMode     string
	MaxMessageSize int64
	MaxRows        int

	PGSQL struct {
		DBName   string
		Username string
		Password string
		Host     string
		Port     string
	}

	Redis struct {
		Host     string
		Port     string
		Username string
		Password string
		DB       int
	}

	Relay struct {
		AdminPubKey string
		Nips        []int

		Name        string
		Description string
		Version     string
		Contract    string
		Software    string
	}
}
