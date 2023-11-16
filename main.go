package main


import (
	"fmt"
	"context"
	"github.com/fvdime/go-study/app"
)

func main(){
	app := app.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}