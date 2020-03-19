package config

import (
	"time"
)

// Config provides the system configuration.
type Config struct {
	Service     Service
	ProxyClient ProxyClient
	Broker      Broker
	Controller  Controller
	//Database Database
	//Logging  Logging
	//Web      Web
}

type Controller struct {
	Period          time.Duration `default:"30s"`
	NewProxyAddress string        `default:"proxy.new"`
}

type Service struct {
	Name        string
	Env         string
	TestLoadEnv bool `envconfig:"TEST_LOAD_ENV"`
	//Logging  Logging
	//Web      Web
}

type ProxyClient struct {
	Url     string        `default:"https://www.sslproxies.org/"`
	Timeout time.Duration `default:"5s"`
}
type Broker struct {
	Url            string        `envconfig:"NATS_URL"`
	AllowReconnect bool          `envconfig:"NATS_ALLOW_RECONNECT"`
	MaxReconnect   int           `envconfig:"NATS_MAX_RECONNECT"`
	ReconnectWait  time.Duration `envconfig:"NATS_RECONNECT_WAIT"`
	Timeout        time.Duration `envconfig:"NATS_TIMEOUT"`
}

//type (
//
//	//// Web provides api server configuration.
//	//Web struct {
//	//	APIHost      string
//	//	ReadTimeout  time.Duration
//	//	WriteTimeout time.Duration
//	//}
//	//// Logging provides the logging configuration.
//	//Logging struct {
//	//	Level  string `envconfig:"LOG_LEVEL"`
//	//	Trace  bool   `envconfig:"DRONE_LOGS_TRACE"`
//	//	Color  bool   `envconfig:"DRONE_LOGS_COLOR"`
//	//	Pretty bool   `envconfig:"DRONE_LOGS_PRETTY"`
//	//	Text   bool   `envconfig:"DRONE_LOGS_TEXT"`
//	//}
//	//// Database provides the database configuration.
//	//Database struct {
//	//	Host       string `envconfig:"PGHOST"`
//	//	Port       string `envconfig:"PGPORT"`
//	//	User       string `envconfig:"PGUSER"`
//	//	Password   string `envconfig:"PGPASSWORD"`
//	//	Database   string `envconfig:"PGDATABASE"`
//	//	AppName    string `envconfig:"PGAPPNAME"`
//	//	DisableTLS bool   `envconfig:"PGDISABLETLS"`
//	//}
//)
//type (
//	// Config provides the system configuration.
//	Config struct {
//		Database Database
//		Logging  Logging
//		Web      Web
//	}
//	// Web provides api server configuration.
//	Web struct {
//		APIHost      string
//		ReadTimeout  time.Duration
//		WriteTimeout time.Duration
//	}
//	// Logging provides the logging configuration.
//	Logging struct {
//		Level  string `envconfig:"LOG_LEVEL"`
//		Trace  bool   `envconfig:"DRONE_LOGS_TRACE"`
//		Color  bool   `envconfig:"DRONE_LOGS_COLOR"`
//		Pretty bool   `envconfig:"DRONE_LOGS_PRETTY"`
//		Text   bool   `envconfig:"DRONE_LOGS_TEXT"`
//	}
//	// Database provides the database configuration.
//	Database struct {
//		Host       string `envconfig:"PGHOST"`
//		Port       string `envconfig:"PGPORT"`
//		User       string `envconfig:"PGUSER"`
//		Password   string `envconfig:"PGPASSWORD"`
//		Database   string `envconfig:"PGDATABASE"`
//		AppName    string `envconfig:"PGAPPNAME"`
//		DisableTLS bool   `envconfig:"PGDISABLETLS"`
//	}
//)
