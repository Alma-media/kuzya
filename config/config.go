package config

type Config struct {
	Broker  Broker
	Storage Storage
}

type Broker struct {
	DSN      string `default:"tcp://localhost:1883"`
	ClientID string `default:"kuzya"`
}


type Storage struct {
	Type     string `default:"memory"`
	Database Database
}

type Database struct {
	Driver string `default:"sqlite3"`
	DSN    string `default:"state.db"`
}
