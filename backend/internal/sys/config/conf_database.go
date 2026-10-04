package config

const (
	DriverSqlite3 = "sqlite3"
)

type Storage struct {
	// Data is the path to the root directory
	Data      string `yaml:"data"       conf:"default:./.data"`
	SqliteURL string `yaml:"sqlite-url" conf:"default:./.data/catalog.db?_pragma=busy_timeout=999&_pragma=journal_mode=WAL&_fk=1"`
}
