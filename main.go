package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alma-media/kuzya/api"
	"github.com/Alma-media/kuzya/config"
	handler "github.com/Alma-media/kuzya/http"
	st "github.com/Alma-media/kuzya/state"
	"github.com/Alma-media/kuzya/state/database"
	"github.com/Alma-media/kuzya/state/database/sqlite"
	"github.com/Alma-media/kuzya/state/memory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	conf "github.com/tiny-go/config"
)

type stateManager interface {
	Switch(deviceID string) (payload string, err error)
	Status(deviceID string) (payload string, err error)
	List() ([]st.Device, error)
}

type endpoint struct {
	name    string
	handler func(string) (string, error)
}

func main() {
	var (
		appConfig config.Config
		state     stateManager
	)

	if err := conf.Init(&appConfig, ""); err != nil {
		log.Fatalf("unable to parse the config: %s", err)
	}

	var db *sql.DB

	switch appConfig.Storage.Type {
	case "memory":
		state = memory.NewSwitch()
	case "database":
		var err error
		db, err = sql.Open(appConfig.Storage.Database.Driver, appConfig.Storage.Database.DSN)
		if err != nil {
			log.Fatalf("unable to establish database connection: %s", err)
		}

		if err := sqlite.Init(context.Background(), db); err != nil {
			log.Fatalf("database migration failure: %s", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("error closing database: %s", err)
			}
		}()

		state = database.NewSwitch(db, new(sqlite.StateManager))
	default:
		log.Fatalf("unknown storage type %q", appConfig.Storage.Type)
	}

	opts := mqtt.
		NewClientOptions().
		AddBroker(appConfig.Broker.DSN).
		SetClientID(appConfig.Broker.ClientID)

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to initialize a client: %s", token.Error())
	}
	defer client.Disconnect(250)

	endpoints := []endpoint{
		{"switch", state.Switch},
		{"status", state.Status},
	}

	for _, endpoint := range endpoints {
		if err := api.RegisterHandler(client, endpoint.name, endpoint.handler); err != nil {
			log.Fatalf("failed to register handler for endpoint %q: %s", endpoint.name, err)
		}
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.HTTP.Port),
		Handler: http.HandlerFunc(handler.NewHandler(state).StatusHandler),
	}

	go func() {
		log.Printf("Started API on port %d", appConfig.HTTP.Port)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	httpServer.Shutdown(shutdownCtx)
}
