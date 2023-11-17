package main


import (
	"fmt"
	"context"
	"github.com/fvdime/go-study/app"
	"os/signal"
	"os"
)

func main(){
	app := app.New(app.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}