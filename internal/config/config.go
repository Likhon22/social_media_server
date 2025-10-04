package config

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type AppConfig struct {
	Addr    string
	DB      *DbConfig
	Version string
	Env     string
}
