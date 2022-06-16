package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"l0/order"
	orderAPI "l0/order/api"
	orderStore "l0/order/repository"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	logger.Info("reading config")
	config, err := NewConfig()
	if err != nil {
		logger.Error("can't decode config", zap.Error(err))
		return
	}

	logger.Info("connecting to database")
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Error("can't open database connection", zap.Error(err), zap.String("db driver", config.DBDriver), zap.String("db source", config.DBSource))
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Error("can't ping database", zap.Error(err), zap.String("db driver", config.DBDriver), zap.String("db source", config.DBSource))
		return
	}

	store := orderStore.New(db)

	sc, err := stan.Connect("test-cluster", "client-456", stan.NatsURL("nats://127.0.0.1:4222"), stan.MaxPubAcksInflight(1000))
	if err != nil {
		logger.Fatal("cat't connect to stan", zap.Error(err))
	}

	sub, err := sc.Subscribe("orders", func(msg *stan.Msg) {
		err := insertMessage(msg.Data, store)
		if err != nil {
			logger.Info("can't create order", zap.Error(err))
		}
		msg.Ack()
	}, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(time.Second))

	if err != nil {
		logger.Fatal("cat't subscribe to channel", zap.Error(err))
	}

	// init router
	api := orderAPI.API{}
	router := api.NewRouter(store)

	// init http server
	srv := &http.Server{
		Addr:        config.HTTPServerAddress,
		Handler:     router,
		ReadTimeout: time.Duration(config.ReadTimeout) * time.Second,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
	}

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("can't start server", zap.Error(err), zap.String("server address", config.HTTPServerAddress))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("received an interrupt, unsubscribing and closing connection")
	sub.Unsubscribe()
	sc.Close()
	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(config.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(timeout); err != nil {
		logger.Error("can't shutdown http server", zap.Error(err))
	}
}

func insertMessage(data []byte, store *orderStore.Queries) error {
	o := new(order.Order)
	err := json.Unmarshal(data, o)
	if err != nil {
		return err
	}

	err = store.CreateOrder(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}
