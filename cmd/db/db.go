package db

type ConnectionConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
}
