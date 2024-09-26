package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/orders-api/application"
)


func main() {
	app := application.New(application.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	err := app.Start(ctx)
	if err != nil {
		fmt.Printf("Failed to listen and serve %v", err)
	}

	defer cancel()
}

// func main() {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger) // logs request calls made to the port.

// 	// using server method
// 	// server := &http.Server{
// 	// 	Addr: ":3000",
// 	// 	Handler: r,
// 	// }
// 	// err := server.ListenAndServe()
// 	// if err != nil {
// 	// 	fmt.Println("Failed and Listen to Server", err)
// 	// }

// 	// using http listen and serve
// 	r.Get("/", basicHandler)
// 	http.ListenAndServe(":3000", r)
// }

// func basicHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello, World"))
// }

// func main() {
// 	server := &http.Server{
// 		Addr: ":3000",
// 		Handler: http.HandlerFunc(basicHandler),
// 	}
// 	err := server.ListenAndServe()
// 	if err != nil {
// 		fmt.Println("Failed to listen to server", err)
// 	}
// }
