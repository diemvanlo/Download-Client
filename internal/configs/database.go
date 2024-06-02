package configs

type DatabaseType string

const (
	DatabaseTypeMySQL DatabaseType = "mysql"
)

type Database struct {
	Type     DatabaseType
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
