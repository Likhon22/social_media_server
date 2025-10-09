package config

import (
	"time"

	"github.com/likhon22/social/internal/mailer"
)

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
	Mailer  mailer.Client
}

type MailConfig struct {
	Exp           time.Duration
	MailGunConfig mailer.MailGunMailer
}
