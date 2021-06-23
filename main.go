package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alma-media/kuzya/api"
	"github.com/Alma-media/kuzya/state/database"
	"github.com/Alma-media/kuzya/state/database/sqlite"
	"github.com/Alma-media/kuzya/state/memory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type stateManager interface {
	Switch(deviceID string) (payload string, err error)
	Status(deviceID string) (payload string, err error)
}

type endpoint struct {
	name    string
	handler func(string) (string, error)
}

func main() {
	var (
		state     stateManager = memory.NewSwitch()
		endpoints              = []endpoint{
			{"switch", state.Switch},
			{"status", state.Status},
		}
	)

	// TODO: parse config variables
	if false {
		db, err := sql.Open("sqlite3", "state.db")
		if err != nil {
			log.Fatalf("unable to establish database connection: %s", err)
		}

		if err := sqlite.Init(context.Background(), db); err != nil {
			log.Fatalf("database migration failure: %s", err)
		}

		defer db.Close()

		state = database.NewSwitch(db, new(sqlite.StateManager))
	}

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("kuzya")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to initialize a client: %s", token.Error())
	}

	for _, endpoint := range endpoints {
		if err := api.RegisterHandler(client, endpoint.name, endpoint.handler); err != nil {
			log.Fatalf("failed to register handler for endpoint %q: %s", endpoint.name, err)
		}
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}
