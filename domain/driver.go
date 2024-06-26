package domain

type Driver struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Db       string
	Ssl      bool
}
