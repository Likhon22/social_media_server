package config

import "time"

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
	ApiURL  string
	Mail    *MailConfig
}

type MailConfig struct {
	Exp time.Duration
}
