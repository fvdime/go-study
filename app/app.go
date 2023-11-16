package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler

	rdb *redis.Client
}

func New() *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{}),
	}

	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3003",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("starting server")

	channel := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			channel <- fmt.Errorf("failed to start the server: %w", err)
		}
		close(channel)
	}()

	select{
	case err = <-channel:
		return err

	case <-	ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

	return nil
}