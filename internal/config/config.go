package config

type Config struct {
	DatabaseDSN string
	LogLevel    string
	JwtSecret   string
	BCryptSalt  string
	AppPort     int
}

func NewConfig() *Config {

}
