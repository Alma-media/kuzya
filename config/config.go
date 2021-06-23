package config

type Config struct {
	Storage Storage
}

type Storage struct {
	Type     string `default:"memory"`
	Database Database
}

type Database struct {
	Driver string `default:"sqlite3"`
	DSN    string `default:"state.db"`
}
