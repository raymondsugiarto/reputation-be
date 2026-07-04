package config

// DatabaseList is struct for list of database
type DatabaseList struct {
	Main  Database
	Redis Database
}

// Database is struct for Database conf
type Database struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Dbname        string
	Schema        string
	Adapter       string
	LogLevel      int
	SlowThreshold int
}
