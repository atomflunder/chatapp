package models

type Config struct {
	Host string
	Port string
}

func GetConfig() Config {
	return Config{
		Host: "http://localhost",
		Port: "8080",
	}
}
