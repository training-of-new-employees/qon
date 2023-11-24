package config

type config struct {
	Address string
}

func New() config {
	return config{
		Address: "127.0.0.1:8080",
	}
}
