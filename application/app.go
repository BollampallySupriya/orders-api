package application

import (
	"net/http"
	"context"
	"fmt"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRouter(),
	}

	return app
}


func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error occurred %w", err)
	}
	return nil
}