package config

type Config struct {
	HTTP    HTTP
	Broker  Broker
	Storage Storage
}

type HTTP struct {
	Port int `default:"8080"`
}

type Broker struct {
	DSN      string `default:"tcp://localhost:1883"`
	ClientID string `default:"kuzya-state-manager"`
}


type Storage struct {
	Type     string `default:"memory"`
	Database Database
}

type Database struct {
	Driver string `default:"sqlite3"`
	DSN    string `default:"state.db"`
}
