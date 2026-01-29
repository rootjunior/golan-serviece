package main

type Config struct {
	ServerAddress string
}

func LoadConfig() Config {
	return Config{
		ServerAddress: ":8080",
	}
}
