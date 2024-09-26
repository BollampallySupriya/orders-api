package application

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
		router: loadRouter(),
		rdb: redis.NewClient(&redis.Options{}),
	}

	return app
}


func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

	redis_err := a.rdb.Ping(ctx).Err()
	if redis_err != nil {
		return fmt.Errorf("error occurred while ping redis server %w", redis_err)
	} else {
		fmt.Println("Started server")
	}

	defer func() {
		if err:= a.rdb.Close(); err != nil {
			fmt.Println("Failed to close redis", err)
		}
	}()

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("error occurred %w", err)
		}
		defer close(ch)
	}()

	// err := server.ListenAndServe()
	// if err != nil {
	// 	return fmt.Errorf("error occurred %w", err)
	// }
	// return nil

	// channel_err:= <-ch // we can use channel_err, open to catch the status of channel open or closed

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second * 10)
		defer cancel()
		return server.Shutdown(timeout)
	}

}
