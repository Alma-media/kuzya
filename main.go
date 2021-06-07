package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/Alma-media/kuzya/state/database"
	"github.com/Alma-media/kuzya/state/database/sqlite"
	"github.com/Alma-media/kuzya/state/memory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var expr = regexp.MustCompile(`^\/trig-in\/(.*)$`)

func createHandler(trig func(string) (string, error)) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if !expr.MatchString(msg.Topic()) {
			log.Println("invalid device format")

			return
		}

		deviceID := expr.FindStringSubmatch(msg.Topic())[1]

		state, err := trig(deviceID)
		if err != nil {
			log.Printf("cannot retrieve current state: %s", err)

			return
		}

		client.Publish("/trig-out/"+deviceID, 0, false, state).Wait()
	}
}

func main() {
	stateSwitch := memory.NewSwitch().Switch

	if true {
		db, err := sql.Open("sqlite3", "state.db")
		if err != nil {
			log.Fatalf("unable to establish database connection: %s", err)
		}

		if err := sqlite.Init(context.Background(), db); err != nil {
			log.Fatalf("database migration failure: %s", err)
		}

		defer db.Close()

		stateSwitch = database.NewSwitch(db, new(sqlite.StateManager)).Switch
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

	token := client.Subscribe("/trig-in/+", 0, createHandler(stateSwitch))

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to subscribe: %s", token.Error())
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}
