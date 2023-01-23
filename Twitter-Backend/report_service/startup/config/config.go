package config

import "os"

type Config struct {
	Port                       string
	ReportDBHost               string
	ReportDBPort               string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	JaegerAddress              string
	CreateReportCommandSubject string
	CreateReportReplySubject   string
}

func NewConfig() *Config {
	return &Config{
		Port:                       os.Getenv("REPORT_SERVICE_PORT"),
		ReportDBHost:               os.Getenv("REPORT_DB_HOST"),
		ReportDBPort:               os.Getenv("REPORT_DB_PORT"),
		NatsHost:                   os.Getenv("NATS_HOST"),
		NatsPort:                   os.Getenv("NATS_PORT"),
		NatsUser:                   os.Getenv("NATS_USER"),
		NatsPass:                   os.Getenv("NATS_PASS"),
		JaegerAddress:              os.Getenv("JAEGER_ADDRESS"),
		CreateReportCommandSubject: os.Getenv("CREATE_REPORT_COMMAND_SUBJECT"),
		CreateReportReplySubject:   os.Getenv("CREATE_REPORT_REPLY_SUBJECT"),
	}
}
