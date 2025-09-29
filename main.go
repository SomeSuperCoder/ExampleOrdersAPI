package main

import (
	"context"
	"log"

	"github.com/SomeSuperCoder/OrdersAPI/application"
)

func main() {
	err := application.New().Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
